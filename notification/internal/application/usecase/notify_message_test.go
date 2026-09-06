package usecase_test

import (
	"chatterbox/notification/internal/application/usecase"
	"chatterbox/notification/internal/domain/entity"
	"chatterbox/notification/internal/domain/port/mocks"
	"chatterbox/notification/internal/domain/valueobject"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNotifyMessageUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		receiverID := valueobject.NewUserID()
		occurredAt := time.Now()

		cmd := usecase.NotifyMessageCommand{
			MessageID:   "message-123",
			ChatID:      "chat-123",
			SenderID:    "sender-123",
			ReceiverIDs: []string{receiverID.String()},
			Text:        "Hello!",
			OccurredAt:  occurredAt,
		}

		sender := mocks.NewMockSender(t)

		sender.EXPECT().
			Send(ctx, mock.MatchedBy(func(n entity.Notification) bool {
				payload, ok := n.Payload.(entity.MessagePayload)
				if !ok {
					return false
				}

				return n.RecepientID == receiverID &&
					n.Type == valueobject.NewMessageNotificationType &&
					payload.ID == cmd.MessageID &&
					payload.ChatID == cmd.ChatID &&
					payload.SenderID == cmd.SenderID &&
					payload.Text == cmd.Text &&
					payload.OccurredAt.Equal(cmd.OccurredAt)
			})).
			Return(nil)

		uc := usecase.NewNotifyMessageUseCase(sender)

		err := uc.Execute(ctx, cmd)

		require.NoError(t, err)
	})

	t.Run("multiple receivers", func(t *testing.T) {
		receiverID1 := valueobject.NewUserID()
		receiverID2 := valueobject.NewUserID()

		cmd := usecase.NotifyMessageCommand{
			MessageID: "message-123",
			ChatID:    "chat-123",
			SenderID:  "sender-123",
			ReceiverIDs: []string{
				receiverID1.String(),
				receiverID2.String(),
			},
			Text:       "Hello!",
			OccurredAt: time.Now(),
		}

		sender := mocks.NewMockSender(t)

		sender.EXPECT().
			Send(ctx, mock.MatchedBy(func(n entity.Notification) bool {
				return n.RecepientID == receiverID1
			})).
			Return(nil)

		sender.EXPECT().
			Send(ctx, mock.MatchedBy(func(n entity.Notification) bool {
				return n.RecepientID == receiverID2
			})).
			Return(nil)

		uc := usecase.NewNotifyMessageUseCase(sender)

		err := uc.Execute(ctx, cmd)

		require.NoError(t, err)
	})

	t.Run("invalid receiver id", func(t *testing.T) {
		cmd := usecase.NotifyMessageCommand{
			ReceiverIDs: []string{"invalid-user-id"},
		}

		sender := mocks.NewMockSender(t)

		uc := usecase.NewNotifyMessageUseCase(sender)

		err := uc.Execute(ctx, cmd)

		assert.Error(t, err)
	})

	t.Run("sender error", func(t *testing.T) {
		receiverID := valueobject.NewUserID()
		expectedErr := errors.New("send error")

		cmd := usecase.NotifyMessageCommand{
			ReceiverIDs: []string{receiverID.String()},
		}

		sender := mocks.NewMockSender(t)

		sender.EXPECT().
			Send(ctx, mock.Anything).
			Return(expectedErr)

		uc := usecase.NewNotifyMessageUseCase(sender)

		err := uc.Execute(ctx, cmd)

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("sender error stops processing", func(t *testing.T) {
		receiverID1 := valueobject.NewUserID()
		receiverID2 := valueobject.NewUserID()
		expectedErr := errors.New("send error")

		cmd := usecase.NotifyMessageCommand{
			ReceiverIDs: []string{
				receiverID1.String(),
				receiverID2.String(),
			},
		}

		sender := mocks.NewMockSender(t)

		sender.EXPECT().
			Send(ctx, mock.Anything).
			Return(expectedErr).
			Once()

		uc := usecase.NewNotifyMessageUseCase(sender)

		err := uc.Execute(ctx, cmd)

		assert.ErrorIs(t, err, expectedErr)
	})
}
