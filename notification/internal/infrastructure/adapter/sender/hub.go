package sender

import "context"

type Hub interface {
	Send(ctx context.Context, userID string, payload interface{}) error
}
