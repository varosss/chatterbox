package port

import "context"

type Event interface {
	Type() string
}

type EventHandler interface {
	Handle(ctx context.Context, e Event) error
}
