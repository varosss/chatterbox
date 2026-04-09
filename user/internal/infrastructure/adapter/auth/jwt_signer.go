package auth

import (
	"chatterbox/user/internal/domain/valueobject"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner struct {
	privateKey *rsa.PrivateKey
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTSigner(
	privateKey *rsa.PrivateKey,
	issuer string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *JWTSigner {
	return &JWTSigner{
		privateKey: privateKey,
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *JWTSigner) SignAccess(
	userID valueobject.UserID,
	now time.Time,
) (string, error) {
	claims := jwt.MapClaims{
		"sub": string(userID),
		"iat": now.Unix(),
		"exp": now.Add(s.accessTTL).Unix(),
		"iss": s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *JWTSigner) SignRefresh(
	tokenID valueobject.TokenID,
	userID valueobject.UserID,
	now time.Time,
) (string, error) {
	claims := jwt.MapClaims{
		"jti": tokenID.String(),
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(s.refreshTTL).Unix(),
		"iss": s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(s.privateKey)
}
