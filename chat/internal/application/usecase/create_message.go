package usecase

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/port"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type CreateMessageCommand struct {
	SenderID valueobject.UserID
	ChatID   valueobject.ChatID
	Text     string
}

type CreateMessageResult struct {
	MessageID valueobject.MessageID
}

type CreateMessageUseCase struct {
	events   port.EventProducer
	chats    port.ChatRepo
	messages port.MessageRepo
}

func NewCreateMessageUseCase(
	events port.EventProducer,
	messages port.MessageRepo,
	chats port.ChatRepo,
) *CreateMessageUseCase {
	return &CreateMessageUseCase{
		events:   events,
		chats:    chats,
		messages: messages,
	}
}

func (uc *CreateMessageUseCase) Execute(
	ctx context.Context,
	cmd CreateMessageCommand,
) (*CreateMessageResult, error) {
	chat, err := uc.chats.FindByID(ctx, cmd.ChatID)
	if err != nil {
		return nil, err
	}

	var receiverIDs []valueobject.UserID
	for _, participantID := range chat.ParticipantIDs() {
		if participantID != cmd.SenderID {
			receiverIDs = append(receiverIDs, participantID)
		}
	}

	message, err := entity.NewMessage(
		cmd.SenderID,
		cmd.ChatID,
		receiverIDs,
		cmd.Text,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messages.Save(ctx, message); err != nil {
		return nil, err
	}

	if err := uc.events.Produce(ctx, message.PullEvents()...); err != nil {
		return nil, err
	}

	return &CreateMessageResult{
		MessageID: message.ID(),
	}, nil
}
