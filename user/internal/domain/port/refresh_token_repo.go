package port

import (
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/valueobject"
	"context"
)

type RefreshTokenRepo interface {
	Save(ctx context.Context, token *entity.RefreshToken) error
	Get(ctx context.Context, id valueobject.TokenID) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id valueobject.TokenID) error
}
