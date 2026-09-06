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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateMessageUseCase_Execute(t *testing.T) {
	t.Run("chat repository error", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(nil, repoErr)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "hello",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("sender is not in chat", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()

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
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "hello",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "sender is not in this chat")
		assert.Nil(t, result)
	})

	t.Run("empty text", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()
		receiverID := valueobject.NewUserID()

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{senderID, receiverID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "text cannot be empty")
		assert.Nil(t, result)
	})

	t.Run("message repository error", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()
		receiverID := valueobject.NewUserID()
		repoErr := errors.New("database error")

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{senderID, receiverID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		messages.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Message")).
			Return(repoErr)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "hello",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("event producer error", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()
		receiverID := valueobject.NewUserID()
		produceErr := errors.New("event producer error")

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{senderID, receiverID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		messages.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Message")).
			Return(nil)

		events.EXPECT().
			Produce(ctx, mock.Anything).
			Return(produceErr)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "hello",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, produceErr)
		assert.Nil(t, result)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		chatID := valueobject.NewChatID()
		senderID := valueobject.NewUserID()
		receiverID := valueobject.NewUserID()

		chat := entity.ChatFromPrimitives(
			chatID,
			[]valueobject.UserID{senderID, receiverID},
			"Test chat",
		)

		chats := mocks.NewMockChatRepo(t)
		messages := mocks.NewMockMessageRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			FindByID(ctx, chatID).
			Return(chat, nil)

		messages.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Message")).
			Return(nil)

		events.EXPECT().
			Produce(ctx, mock.Anything).
			Return(nil)

		uc := usecase.NewCreateMessageUseCase(events, messages, chats)

		result, err := uc.Execute(ctx, usecase.CreateMessageCommand{
			SenderID: senderID,
			ChatID:   chatID,
			Text:     "hello",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.MessageID.String())
	})
}
