package usecase

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/port"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type CreateChatCommand struct {
	ParticipantIDs []valueobject.UserID
	DisplayName    string
}

type CreateChatResult struct {
	ChatID valueobject.ChatID
}

type CreateChatUseCase struct {
	events port.EventProducer
	chats  port.ChatRepo
}

func NewCreateChatUseCase(
	events port.EventProducer,
	chats port.ChatRepo,
) *CreateChatUseCase {
	return &CreateChatUseCase{
		events: events,
		chats:  chats,
	}
}

func (uc *CreateChatUseCase) Execute(
	ctx context.Context,
	cmd CreateChatCommand,
) (*CreateChatResult, error) {
	chat, err := entity.NewChat(
		cmd.ParticipantIDs,
		cmd.DisplayName,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.chats.Save(ctx, chat); err != nil {
		return nil, err
	}

	if err := uc.events.Produce(ctx, chat.PullEvents()...); err != nil {
		return nil, err
	}

	return &CreateChatResult{
		ChatID: chat.ID(),
	}, nil
}
