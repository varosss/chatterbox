package port

import (
	"chatterbox/notification/internal/domain/entity"
	"context"
)

type Sender interface {
	Send(ctx context.Context, n entity.Notification) error
}
