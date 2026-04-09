package auth

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	publicKey *rsa.PublicKey
	issuer    string
}

func NewJWTVerifier(
	publicKey *rsa.PublicKey,
	issuer string,
) *JWTVerifier {
	return &JWTVerifier{
		publicKey: publicKey,
		issuer:    issuer,
	}
}

func (v *JWTVerifier) VerifyAccess(
	tokenString string,
) (*AccessTokenClaims, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("invalid token")
			}
			return v.publicKey, nil
		},
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, errors.New("invalid claims")
	}

	return &AccessTokenClaims{
		UserID: sub,
	}, nil
}
