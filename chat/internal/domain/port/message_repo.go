package port

import (
	"chatterbox/chat/internal/domain/entity"
	"context"
)

type MessageRepo interface {
	Save(ctx context.Context, msg *entity.Message) error
}
