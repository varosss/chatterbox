package usecase_test

import (
	"context"
	"errors"
	"testing"

	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/port/mocks"
	"chatterbox/chat/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListChatsUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		userID := valueobject.NewUserID()
		participantID := valueobject.NewUserID()
		chatID := valueobject.NewChatID()

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{userID, participantID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)

		chats.EXPECT().
			List(ctx, userID).
			Return([]*entity.Chat{chat}, nil)

		uc := usecase.NewListChatsUseCase(chats)

		result, err := uc.Execute(ctx, usecase.ListChatsCommand{
			UserID: userID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Chats, 1)

		assert.Equal(t, chatID.String(), result.Chats[0].ID)
		assert.Equal(t, "Test chat", result.Chats[0].DisplayName)
		assert.Equal(t,
			[]string{userID.String(), participantID.String()},
			result.Chats[0].ParticipantIDs,
		)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()

		userID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		chats := mocks.NewMockChatRepo(t)

		chats.EXPECT().
			List(ctx, userID).
			Return(nil, repoErr)

		uc := usecase.NewListChatsUseCase(chats)

		result, err := uc.Execute(ctx, usecase.ListChatsCommand{
			UserID: userID,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})
}
