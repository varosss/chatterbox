package usecase

import (
	"chatterbox/notification/internal/domain/entity"
	"chatterbox/notification/internal/domain/port"
	"chatterbox/notification/internal/domain/valueobject"
	"context"
	"time"
)

type NotifyMessageCommand struct {
	MessageID   string
	ChatID      string
	SenderID    string
	ReceiverIDs []string
	Text        string
	OccurredAt  time.Time
}

type NotifyMessageUseCase struct {
	sender port.Sender
}

func NewNotifyMessageUseCase(sender port.Sender) *NotifyMessageUseCase {
	return &NotifyMessageUseCase{
		sender: sender,
	}
}

func (uc *NotifyMessageUseCase) Execute(
	ctx context.Context,
	cmd NotifyMessageCommand,
) error {
	for _, userID := range cmd.ReceiverIDs {
		recepientID, err := valueobject.ParseUserID(userID)
		if err != nil {
			return err
		}

		n := entity.Notification{
			ID:          valueobject.NewNotificationID(),
			RecepientID: recepientID,
			Type:        valueobject.NewMessageNotificationType,
			Payload: entity.MessagePayload{
				ChatID:     cmd.ChatID,
				ID:         cmd.MessageID,
				SenderID:   cmd.SenderID,
				Text:       cmd.Text,
				OccurredAt: cmd.OccurredAt,
			},
		}

		if err := uc.sender.Send(ctx, n); err != nil {
			return err
		}
	}

	return nil
}
