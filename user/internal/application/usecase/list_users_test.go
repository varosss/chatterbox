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

func TestListUsersUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		userID1 := valueobject.NewUserID()
		userID2 := valueobject.NewUserID()

		email1, err := valueobject.NewEmail("john@example.com")
		require.NoError(t, err)

		email2, err := valueobject.NewEmail("jane@example.com")
		require.NoError(t, err)

		passwordHash, err := valueobject.NewPasswordHash("hashed-password")
		require.NoError(t, err)

		user1 := entity.UserFromPrimitives(
			userID1,
			email1,
			"john",
			"John Doe",
			passwordHash,
			valueobject.ActiveStatus,
		)

		user2 := entity.UserFromPrimitives(
			userID2,
			email2,
			"jane",
			"Jane Doe",
			passwordHash,
			valueobject.BlockedStatus,
		)

		userIDs := []valueobject.UserID{userID1, userID2}

		users := mocks.NewMockUserRepo(t)

		users.EXPECT().
			List(ctx, userIDs, 100, 20).
			Return([]*entity.User{user1, user2}, nil)

		uc := usecase.NewListUsersUseCase(users)

		result, err := uc.Execute(ctx, usecase.ListUsersCommand{
			UserIDs: userIDs,
			Limit:   100,
			Offset:  20,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Users, 2)

		assert.Equal(t, userID1.String(), result.Users[0].ID)
		assert.Equal(t, email1.String(), result.Users[0].Email)
		assert.Equal(t, "john", result.Users[0].Username)
		assert.Equal(t, "John Doe", result.Users[0].DisplayName)
		assert.Equal(t, valueobject.ActiveStatus.Int(), result.Users[0].Status)

		assert.Equal(t, userID2.String(), result.Users[1].ID)
		assert.Equal(t, email2.String(), result.Users[1].Email)
		assert.Equal(t, "jane", result.Users[1].Username)
		assert.Equal(t, "Jane Doe", result.Users[1].DisplayName)
		assert.Equal(t, valueobject.BlockedStatus.Int(), result.Users[1].Status)
	})

	t.Run("default limit", func(t *testing.T) {
		ctx := context.Background()

		users := mocks.NewMockUserRepo(t)

		var userIDs []valueobject.UserID = nil

		users.EXPECT().
			List(ctx, userIDs, 500, 0).
			Return([]*entity.User{}, nil)

		uc := usecase.NewListUsersUseCase(users)

		result, err := uc.Execute(ctx, usecase.ListUsersCommand{})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Users)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()

		repoErr := errors.New("database error")

		users := mocks.NewMockUserRepo(t)

		var userIDs []valueobject.UserID = nil

		users.EXPECT().
			List(ctx, userIDs, 100, 0).
			Return(nil, repoErr)

		uc := usecase.NewListUsersUseCase(users)

		result, err := uc.Execute(ctx, usecase.ListUsersCommand{
			Limit: 100,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})
}
