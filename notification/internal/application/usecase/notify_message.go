package usecase

import (
	"chatterbox/notification/internal/domain/entity"
	"chatterbox/notification/internal/domain/port"
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
	for _, receiverID := range cmd.ReceiverIDs {
		n := entity.Notification{
			UserID: receiverID,
			Type:   "new_message",
			Data: entity.MessageData{
				ChatID:   cmd.ChatID,
				ID:       cmd.MessageID,
				SenderID: cmd.SenderID,
				Text:     cmd.Text,
			},
		}

		if err := uc.sender.Send(ctx, n); err != nil {
			return err
		}
	}

	return nil
}
