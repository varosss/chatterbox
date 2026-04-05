package eventhandler

import (
	"chatterbox/notification/internal/application/port"
	"chatterbox/notification/internal/application/usecase"
	"chatterbox/notification/internal/infrastructure/controller/event"
	"context"
	"fmt"
)

type MessageCreatedHandler struct {
	uc *usecase.NotifyMessageUseCase
}

func NewMessageCreatedHandler(uc *usecase.NotifyMessageUseCase) *MessageCreatedHandler {
	return &MessageCreatedHandler{
		uc: uc,
	}
}

func (h *MessageCreatedHandler) Handle(ctx context.Context, e port.Event) error {
	messageCreated, ok := e.(event.MessageCreated)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	if err := h.uc.Execute(ctx, usecase.NotifyMessageCommand{
		MessageID:   messageCreated.MessageID,
		ChatID:      messageCreated.ChatID,
		SenderID:    messageCreated.SenderID,
		ReceiverIDs: messageCreated.Receivers,
		Text:        messageCreated.Text,
		OccurredAt:  messageCreated.OccurredAt,
	}); err != nil {
		return err
	}

	return nil
}
