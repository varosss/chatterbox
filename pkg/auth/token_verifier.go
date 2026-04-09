package auth

type AccessTokenClaims struct {
	UserID string
}

type TokenVerifier interface {
	VerifyAccess(token string) (*AccessTokenClaims, error)
}
