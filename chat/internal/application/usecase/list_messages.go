package usecase

import (
	"chatterbox/chat/internal/application/dto"
	"chatterbox/chat/internal/domain/port"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type ListMessagesCommand struct {
	ChatID valueobject.ChatID
}

type ListMessagesResult struct {
	Messages []*dto.Message
}

type ListMessagesUseCase struct {
	messages port.MessageRepo
}

func NewListMessagesUseCase(messages port.MessageRepo) *ListMessagesUseCase {
	return &ListMessagesUseCase{
		messages: messages,
	}
}

func (uc *ListMessagesUseCase) Execute(
	ctx context.Context,
	cmd ListMessagesCommand,
) (*ListMessagesResult, error) {
	messages, err := uc.messages.FindManyByChatID(ctx, cmd.ChatID)
	if err != nil {
		return nil, err
	}

	messagesRes := make([]*dto.Message, len(messages))
	for i, message := range messages {
		messagesRes[i] = &dto.Message{
			ID:        message.ID().String(),
			ChatID:    message.ChatID().String(),
			SenderID:  message.SenderID().String(),
			Text:      message.Text(),
			CreatedAt: message.CreatedAt(),
		}
	}

	return &ListMessagesResult{
		Messages: messagesRes,
	}, nil
}
