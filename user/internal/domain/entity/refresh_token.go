package entity

import (
	"chatterbox/user/internal/domain/valueobject"
	"time"
)

type RefreshToken struct {
	id        valueobject.TokenID
	userID    valueobject.UserID
	expiresAt time.Time
	revoked   bool
}

func NewRefreshToken(
	userID valueobject.UserID,
	expiresAt time.Time,
) *RefreshToken {
	return &RefreshToken{
		id:        valueobject.NewTokenID(),
		userID:    userID,
		expiresAt: expiresAt,
	}
}

func RefreshTokenFromPrimitives(
	id valueobject.TokenID,
	userID valueobject.UserID,
	expiresAt time.Time,
	revoked bool,
) *RefreshToken {
	return &RefreshToken{
		id:        id,
		userID:    userID,
		expiresAt: expiresAt,
		revoked:   revoked,
	}
}

func (t *RefreshToken) ID() valueobject.TokenID {
	return t.id
}

func (t *RefreshToken) UserID() valueobject.UserID {
	return t.userID
}

func (t *RefreshToken) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t *RefreshToken) IsRevoked() bool {
	return t.revoked
}

func (t *RefreshToken) IsExpired(now time.Time) bool {
	return now.After(t.expiresAt)
}

func (t *RefreshToken) Revoke() {
	t.revoked = true
}
