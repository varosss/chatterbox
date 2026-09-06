package usecase_test

import (
	"context"
	"errors"
	"testing"

	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/port/mocks"
	"chatterbox/chat/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateChatUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		user1 := valueobject.NewUserID()
		user2 := valueobject.NewUserID()

		chats := mocks.NewMockChatRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Chat")).
			Return(nil)

		events.EXPECT().
			Produce(ctx, mock.Anything).
			Return(nil)

		uc := usecase.NewCreateChatUseCase(events, chats)

		result, err := uc.Execute(ctx, usecase.CreateChatCommand{
			ParticipantIDs: []valueobject.UserID{user1, user2},
			DisplayName:    "Test chat",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.ChatID.String())
	})

	t.Run("cannot create chat without participants", func(t *testing.T) {
		ctx := context.Background()

		chats := mocks.NewMockChatRepo(t)
		events := mocks.NewMockEventProducer(t)

		uc := usecase.NewCreateChatUseCase(events, chats)

		result, err := uc.Execute(ctx, usecase.CreateChatCommand{
			DisplayName: "Test chat",
		})

		require.Error(t, err)
		assert.EqualError(t, err, "cannot create chat without participants")
		assert.Nil(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		repoErr := errors.New("database error")

		chats := mocks.NewMockChatRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Chat")).
			Return(repoErr)

		uc := usecase.NewCreateChatUseCase(events, chats)

		result, err := uc.Execute(ctx, usecase.CreateChatCommand{
			ParticipantIDs: []valueobject.UserID{
				valueobject.NewUserID(),
			},
			DisplayName: "Test chat",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("event producer error", func(t *testing.T) {
		ctx := context.Background()
		produceErr := errors.New("event producer error")

		chats := mocks.NewMockChatRepo(t)
		events := mocks.NewMockEventProducer(t)

		chats.EXPECT().
			Save(ctx, mock.AnythingOfType("*entity.Chat")).
			Return(nil)

		events.EXPECT().
			Produce(ctx, mock.Anything).
			Return(produceErr)

		uc := usecase.NewCreateChatUseCase(events, chats)

		result, err := uc.Execute(ctx, usecase.CreateChatCommand{
			ParticipantIDs: []valueobject.UserID{
				valueobject.NewUserID(),
			},
			DisplayName: "Test chat",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, produceErr)
		assert.Nil(t, result)
	})
}
