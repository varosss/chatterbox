package usecase

import (
	"chatterbox/chat/internal/application/dto"
	"chatterbox/chat/internal/domain/port"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
	"errors"
	"slices"
)

type ListMessagesCommand struct {
	UserID valueobject.UserID
	ChatID valueobject.ChatID
}

type ListMessagesResult struct {
	Messages []*dto.Message
}

type ListMessagesUseCase struct {
	chats    port.ChatRepo
	messages port.MessageRepo
}

func NewListMessagesUseCase(
	chats port.ChatRepo,
	messages port.MessageRepo,
) *ListMessagesUseCase {
	return &ListMessagesUseCase{
		chats:    chats,
		messages: messages,
	}
}

func (uc *ListMessagesUseCase) Execute(
	ctx context.Context,
	cmd ListMessagesCommand,
) (*ListMessagesResult, error) {
	chat, err := uc.chats.FindByID(ctx, cmd.ChatID)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(chat.ParticipantIDs(), cmd.UserID) {
		return nil, errors.New("user is not a chat participant")
	}

	messages, err := uc.messages.List(ctx, chat.ID())
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
