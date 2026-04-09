package port

import "chatterbox/user/internal/domain/valueobject"

type AccessTokenClaims struct {
	UserID valueobject.UserID
}

type TokenVerifier interface {
	VerifyAccess(token string) (*AccessTokenClaims, error)
	VerifyRefresh(token string) (valueobject.TokenID, valueobject.UserID, error)
}
