package port

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
)

type MessageRepo interface {
	Save(ctx context.Context, msg *entity.Message) error
	List(ctx context.Context, chatID valueobject.ChatID) ([]*entity.Message, error)
}
