package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/port/mocks"
	"chatterbox/user/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	tokenID := valueobject.NewTokenID()
	userID := valueobject.NewUserID()
	now := time.Now()
	ttl := 24 * time.Hour

	t.Run("invalid token", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		verifier.EXPECT().
			VerifyRefresh("bad").
			Return(tokenID, userID, errors.New("invalid"))

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "bad",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "invalid token")
		assert.Nil(t, result)
	})

	t.Run("token not found", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		verifier.EXPECT().
			VerifyRefresh("refresh").
			Return(tokenID, userID, nil)

		refreshTokens.EXPECT().
			Get(ctx, tokenID).
			Return(nil, errors.New("not found"))

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "refresh",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "invalid token")
		assert.Nil(t, result)
	})

	t.Run("revoked token", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		token := entity.RefreshTokenFromPrimitives(
			tokenID,
			userID,
			now.Add(time.Hour),
			true,
		)

		verifier.EXPECT().VerifyRefresh("refresh").Return(tokenID, userID, nil)
		refreshTokens.EXPECT().Get(ctx, tokenID).Return(token, nil)
		clock.EXPECT().Now().Return(now)

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "refresh",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "invalid token")
		assert.Nil(t, result)
	})

	t.Run("expired token", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		token := entity.RefreshTokenFromPrimitives(
			tokenID,
			userID,
			now.Add(-time.Hour),
			false,
		)

		verifier.EXPECT().VerifyRefresh("refresh").Return(tokenID, userID, nil)
		refreshTokens.EXPECT().Get(ctx, tokenID).Return(token, nil)
		clock.EXPECT().Now().Return(now)

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "refresh",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "invalid token")
		assert.Nil(t, result)
	})

	t.Run("revoke error", func(t *testing.T) {
		revokeErr := errors.New("database error")

		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		token := entity.RefreshTokenFromPrimitives(
			tokenID,
			userID,
			now.Add(time.Hour),
			false,
		)

		verifier.EXPECT().VerifyRefresh("refresh").Return(tokenID, userID, nil)
		refreshTokens.EXPECT().Get(ctx, tokenID).Return(token, nil)
		clock.EXPECT().Now().Return(now)
		refreshTokens.EXPECT().Revoke(ctx, tokenID).Return(revokeErr)

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "refresh",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, revokeErr)
		assert.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		token := entity.RefreshTokenFromPrimitives(
			tokenID,
			userID,
			now.Add(time.Hour),
			false,
		)

		verifier.EXPECT().VerifyRefresh("refresh").Return(tokenID, userID, nil)
		refreshTokens.EXPECT().Get(ctx, tokenID).Return(token, nil)
		clock.EXPECT().Now().Return(now)

		refreshTokens.EXPECT().Revoke(ctx, tokenID).Return(nil)
		signer.EXPECT().SignAccess(userID, now).Return("access-token", nil)
		refreshTokens.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.RefreshToken")).
			Return(nil)
		signer.EXPECT().
			SignRefresh(mock.Anything, userID, now).
			Return("new-refresh-token", nil)

		uc := usecase.NewRefreshTokenUseCase(
			refreshTokens, verifier, signer, clock, ttl,
		)

		result, err := uc.Execute(ctx, usecase.RefreshTokenCommand{
			RefreshToken: "refresh",
		})

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "access-token", result.AccessToken)
		assert.Equal(t, "new-refresh-token", result.RefreshToken)
		assert.True(t, token.IsRevoked())
	})
}
