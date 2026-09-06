package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/port/mocks"
	"chatterbox/chat/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMessagesUseCase_Execute(t *testing.T) {
	t.Run("chat repository error", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		userID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(nil, repoErr)

		uc := usecase.NewListMessagesUseCase(chats, messages)

		result, err := uc.Execute(ctx, usecase.ListMessagesCommand{
			UserID: userID,
			ChatID: chatID,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("user is not a chat participant", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		userID := valueobject.NewUserID()

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{
				valueobject.NewUserID(),
				valueobject.NewUserID(),
			},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		uc := usecase.NewListMessagesUseCase(chats, messages)

		result, err := uc.Execute(ctx, usecase.ListMessagesCommand{
			UserID: userID,
			ChatID: chatID,
		})

		require.Error(t, err)
		assert.EqualError(t, err, "user is not a chat participant")
		assert.Nil(t, result)
	})

	t.Run("message repository error", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		userID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{userID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		messages.EXPECT().
			List(ctx, chatID).
			Return(nil, repoErr)

		uc := usecase.NewListMessagesUseCase(chats, messages)

		result, err := uc.Execute(ctx, usecase.ListMessagesCommand{
			UserID: userID,
			ChatID: chatID,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		userID := valueobject.NewUserID()
		messageID := valueobject.NewMessageID()

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{userID},
			"Test chat",
		)

		message := entity.MessageFromPrimitives(
			messageID,
			userID,
			chatID,
			"Hello",
			time.Now(),
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		messages.EXPECT().
			List(ctx, chatID).
			Return([]*entity.Message{message}, nil)

		uc := usecase.NewListMessagesUseCase(chats, messages)

		result, err := uc.Execute(ctx, usecase.ListMessagesCommand{
			UserID: userID,
			ChatID: chatID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Messages, 1)

		assert.Equal(t, messageID.String(), result.Messages[0].ID)
		assert.Equal(t, chatID.String(), result.Messages[0].ChatID)
		assert.Equal(t, userID.String(), result.Messages[0].SenderID)
		assert.Equal(t, "Hello", result.Messages[0].Text)
		assert.Equal(t, message.CreatedAt(), result.Messages[0].CreatedAt)
	})
}
