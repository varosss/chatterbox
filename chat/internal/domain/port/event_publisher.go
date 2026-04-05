package port

import (
	"chatterbox/chat/internal/domain/event"
	"context"
)

type EventProducer interface {
	Produce(ctx context.Context, events ...event.Event) error
}
