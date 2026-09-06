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

func TestLoginUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	email, _ := valueobject.NewEmail("test@example.com")
	passwordHash, _ := valueobject.NewPasswordHash("hashed-password")
	userID := valueobject.NewUserID()
	now := time.Now()
	refreshTTL := 24 * time.Hour

	makeUser := func() *entity.User {
		return entity.UserFromPrimitives(
			userID,
			email,
			"john",
			"John Doe",
			passwordHash,
			valueobject.ActiveStatus,
		)
	}

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("database error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, repoErr)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email:    email,
			Password: "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("invalid password", func(t *testing.T) {
		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().
			FindByEmail(ctx, email).
			Return(makeUser(), nil)

		passwords.EXPECT().
			Compare("hashed-password", "password").
			Return(false)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email:    email,
			Password: "password",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "password is invalid")
		assert.Nil(t, result)
	})

	t.Run("sign access error", func(t *testing.T) {
		signErr := errors.New("sign error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().FindByEmail(ctx, email).Return(makeUser(), nil)
		passwords.EXPECT().Compare("hashed-password", "password").Return(true)
		clock.EXPECT().Now().Return(now)
		signer.EXPECT().SignAccess(userID, now).Return("", signErr)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email: email, Password: "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, signErr)
		assert.Nil(t, result)
	})

	t.Run("save refresh token error", func(t *testing.T) {
		saveErr := errors.New("save error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().FindByEmail(ctx, email).Return(makeUser(), nil)
		passwords.EXPECT().Compare("hashed-password", "password").Return(true)
		clock.EXPECT().Now().Return(now)
		signer.EXPECT().SignAccess(userID, now).Return("access", nil)
		refreshTokens.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.RefreshToken")).
			Return(saveErr)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email: email, Password: "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, saveErr)
		assert.Nil(t, result)
	})

	t.Run("sign refresh error", func(t *testing.T) {
		signErr := errors.New("refresh sign error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().FindByEmail(ctx, email).Return(makeUser(), nil)
		passwords.EXPECT().Compare("hashed-password", "password").Return(true)
		clock.EXPECT().Now().Return(now)
		signer.EXPECT().SignAccess(userID, now).Return("access", nil)
		refreshTokens.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.RefreshToken")).
			Return(nil)
		signer.EXPECT().
			SignRefresh(mock.Anything, userID, now).
			Return("", signErr)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email: email, Password: "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, signErr)
		assert.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordVerifier(t)
		refreshTokens := mocks.NewMockRefreshTokenRepo(t)
		signer := mocks.NewMockTokenSigner(t)
		clock := mocks.NewMockClock(t)

		users.EXPECT().FindByEmail(ctx, email).Return(makeUser(), nil)
		passwords.EXPECT().Compare("hashed-password", "password").Return(true)
		clock.EXPECT().Now().Return(now)
		signer.EXPECT().SignAccess(userID, now).Return("access-token", nil)
		refreshTokens.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.RefreshToken")).
			Return(nil)
		signer.EXPECT().
			SignRefresh(mock.Anything, userID, now).
			Return("refresh-token", nil)

		uc := usecase.NewLoginUseCase(
			users, passwords, refreshTokens, signer, clock, refreshTTL,
		)

		result, err := uc.Execute(ctx, usecase.LoginCommand{
			Email: email, Password: "password",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "access-token", result.AccessToken)
		assert.Equal(t, "refresh-token", result.RefreshToken)
	})
}
