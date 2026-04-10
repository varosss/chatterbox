package auth

import (
	httpmiddleware "chatterbox/pkg/auth"
	"chatterbox/user/internal/domain/port"
)

type TokenVerifierWrapper struct {
	verifier port.TokenVerifier
}

func NewTokenVerifierWrapper(
	verifier port.TokenVerifier,
) *TokenVerifierWrapper {
	return &TokenVerifierWrapper{
		verifier: verifier,
	}
}

func (a *TokenVerifierWrapper) VerifyAccess(
	token string,
) (*httpmiddleware.AccessTokenClaims, error) {
	claims, err := a.verifier.VerifyAccess(token)
	if err != nil {
		return nil, err
	}

	return &httpmiddleware.AccessTokenClaims{
		UserID: claims.UserID.String(),
	}, nil
}
