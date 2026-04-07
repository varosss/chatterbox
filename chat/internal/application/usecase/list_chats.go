package usecase

import (
	"chatterbox/chat/internal/application/dto"
	"chatterbox/chat/internal/domain/port"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type ListChatsCommand struct {
	UserID valueobject.UserID
}

type ListChatsResult struct {
	Chats []*dto.Chat
}

type ListChatsUseCase struct {
	chats port.ChatRepo
}

func NewListChatsUseCase(chats port.ChatRepo) *ListChatsUseCase {
	return &ListChatsUseCase{
		chats: chats,
	}
}

func (uc *ListChatsUseCase) Execute(
	ctx context.Context,
	cmd ListChatsCommand,
) (*ListChatsResult, error) {
	chats, err := uc.chats.FindManyByParticipantID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	chatsRes := make([]*dto.Chat, len(chats))
	for i, chat := range chats {
		chatsRes[i] = &dto.Chat{
			ID:             chat.ID().String(),
			ParticipantIDs: chat.ParticipantIDsAsStrings(),
		}
	}

	return &ListChatsResult{
		Chats: chatsRes,
	}, nil
}
