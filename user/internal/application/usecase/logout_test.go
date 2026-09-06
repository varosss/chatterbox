package usecase_test

import (
	"context"
	"errors"
	"testing"

	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/port/mocks"
	"chatterbox/user/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	tokenID := valueobject.NewTokenID()
	userID := valueobject.NewUserID()

	t.Run("invalid token", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)

		verifier.EXPECT().
			VerifyRefresh("bad-token").
			Return(tokenID, userID, errors.New("invalid"))

		uc := usecase.NewLogoutUseCase(verifier, refreshTokens)

		err := uc.Execute(ctx, usecase.LogoutCommand{
			RefreshToken: "bad-token",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "invalid token")
	})

	t.Run("revoke error", func(t *testing.T) {
		revokeErr := errors.New("database error")

		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)

		verifier.EXPECT().
			VerifyRefresh("refresh").
			Return(tokenID, userID, nil)

		refreshTokens.EXPECT().
			Revoke(ctx, tokenID).
			Return(revokeErr)

		uc := usecase.NewLogoutUseCase(verifier, refreshTokens)

		err := uc.Execute(ctx, usecase.LogoutCommand{
			RefreshToken: "refresh",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, revokeErr)
	})

	t.Run("success", func(t *testing.T) {
		verifier := mocks.NewMockTokenVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)

		verifier.EXPECT().
			VerifyRefresh("refresh").
			Return(tokenID, userID, nil)

		refreshTokens.EXPECT().
			Revoke(ctx, tokenID).
			Return(nil)

		uc := usecase.NewLogoutUseCase(verifier, refreshTokens)

		err := uc.Execute(ctx, usecase.LogoutCommand{
			RefreshToken: "refresh",
		})

		require.NoError(t, err)
	})
}
