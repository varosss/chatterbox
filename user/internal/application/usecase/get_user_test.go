package usecase_test

import (
	"context"
	"errors"
	"testing"

	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/port/mocks"
	"chatterbox/user/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		userID := valueobject.NewUserID()
		email, err := valueobject.NewEmail("test@example.com")
		require.NoError(t, err)

		passwordHash, err := valueobject.NewPasswordHash("hashed-password")
		require.NoError(t, err)

		user := entity.UserFromPrimitives(
			userID,
			email,
			"john",
			"John Doe",
			passwordHash,
			valueobject.ActiveStatus,
		)

		users := mocks.NewMockUserRepo(t)

		users.EXPECT().
			FindByID(ctx, userID).
			Return(user, nil)

		uc := usecase.NewGetUserUseCase(users)

		result, err := uc.Execute(ctx, usecase.GetUserCommand{
			UserID: userID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.User)

		assert.Equal(t, userID.String(), result.User.ID)
		assert.Equal(t, email.String(), result.User.Email)
		assert.Equal(t, "john", result.User.Username)
		assert.Equal(t, "John Doe", result.User.DisplayName)
		assert.Equal(t, valueobject.ActiveStatus.Int(), result.User.Status)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()

		userID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		users := mocks.NewMockUserRepo(t)

		users.EXPECT().
			FindByID(ctx, userID).
			Return(nil, repoErr)

		uc := usecase.NewGetUserUseCase(users)

		result, err := uc.Execute(ctx, usecase.GetUserCommand{
			UserID: userID,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})
}
