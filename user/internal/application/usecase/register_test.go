package usecase_test

import (
	"context"
	"errors"
	"testing"

	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/port/mocks"
	"chatterbox/user/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	t.Run("user already exists", func(t *testing.T) {
		ctx := context.Background()

		email, err := valueobject.NewEmail("test@example.com")
		require.NoError(t, err)

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordHasher(t)

		users.EXPECT().
			ExistsByEmail(ctx, email).
			Return(true)

		uc := usecase.NewRegisterUseCase(users, passwords)

		result, err := uc.Execute(ctx, usecase.RegisterCommand{
			Email:    email,
			Username: "john",
			Password: "password",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "user already exists")
		assert.Nil(t, result)
	})

	t.Run("hash error", func(t *testing.T) {
		ctx := context.Background()

		email, err := valueobject.NewEmail("test@example.com")
		require.NoError(t, err)

		hashErr := errors.New("hash error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordHasher(t)

		users.EXPECT().
			ExistsByEmail(ctx, email).
			Return(false)

		passwords.EXPECT().
			Hash("password").
			Return("", hashErr)

		uc := usecase.NewRegisterUseCase(users, passwords)

		result, err := uc.Execute(ctx, usecase.RegisterCommand{
			Email:    email,
			Username: "john",
			Password: "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, hashErr)
		assert.Nil(t, result)
	})

	t.Run("save error", func(t *testing.T) {
		ctx := context.Background()

		email, err := valueobject.NewEmail("test@example.com")
		require.NoError(t, err)

		saveErr := errors.New("database error")

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordHasher(t)

		users.EXPECT().
			ExistsByEmail(ctx, email).
			Return(false)

		passwords.EXPECT().
			Hash("password").
			Return("hashed-password", nil)

		users.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.User")).
			Return(saveErr)

		uc := usecase.NewRegisterUseCase(users, passwords)

		result, err := uc.Execute(ctx, usecase.RegisterCommand{
			Email:       email,
			Username:    "john",
			DisplayName: "John Doe",
			Password:    "password",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, saveErr)
		assert.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		email, err := valueobject.NewEmail("test@example.com")
		require.NoError(t, err)

		users := mocks.NewMockUserRepo(t)
		passwords := mocks.NewMockPasswordHasher(t)

		users.EXPECT().
			ExistsByEmail(ctx, email).
			Return(false)

		passwords.EXPECT().
			Hash("password").
			Return("hashed-password", nil)

		users.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.User")).
			Return(nil)

		uc := usecase.NewRegisterUseCase(users, passwords)

		result, err := uc.Execute(ctx, usecase.RegisterCommand{
			Email:       email,
			Username:    "john",
			DisplayName: "John Doe",
			Password:    "password",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.UserID)
	})
}
