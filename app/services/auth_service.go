package services

import (
	"broker/app/types"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService interface {
	CreateToken(ctx context.Context, id string) (string, error)
}

type AuthServiceImpl struct {
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthService(
	signingKey []byte,
	expireDuration time.Duration,
) AuthService {
	return &AuthServiceImpl{
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *AuthServiceImpl) CreateToken(ctx context.Context, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Id: id,
	})

	return token.SignedString(a.signingKey)
}
