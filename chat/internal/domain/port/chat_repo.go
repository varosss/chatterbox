package port

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type ChatRepo interface {
	Save(ctx context.Context, chat *entity.Chat) error
	FindByID(ctx context.Context, chatID valueobject.ChatID) (*entity.Chat, error)
}
