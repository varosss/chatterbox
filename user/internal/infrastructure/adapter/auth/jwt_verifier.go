package auth

import (
	"chatterbox/user/internal/domain/port"
	"chatterbox/user/internal/domain/valueobject"
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
) (*port.AccessTokenClaims, error) {

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

	return &port.AccessTokenClaims{
		UserID: valueobject.UserID(sub),
	}, nil
}

func (v *JWTVerifier) VerifyRefresh(
	tokenString string,
) (valueobject.TokenID, valueobject.UserID, error) {

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
	if err != nil {
		return "", "", errors.New("invalid token")
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", "", errors.New("invalid claims")
	}

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return "", "", errors.New("invalid claims")
	}

	return valueobject.TokenID(jti), valueobject.UserID(sub), nil
}
