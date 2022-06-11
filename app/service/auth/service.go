package auth

import (
	"broker/app/types"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)


type AuthServiceImpl struct {
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthService(
	signingKey []byte,
	expireDuration time.Duration,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *AuthServiceImpl) Refresh(id string) (map[string]string, error) {
	auth := map[string]string{}
	accessToken, err := a.createToken(id, a.expireDuration)
	if err != nil {
		return auth, err
	}
	auth["accessToken"] = accessToken
	refreshToken, err := a.createToken(id, time.Duration(24*30))
	if err != nil {
		return auth, err
	}
	auth["refreshToken"] = refreshToken
	return auth, nil
}

func (a *AuthServiceImpl) createToken(id string, expireAt time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expireAt)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Id: id,
	})

	return token.SignedString(a.signingKey)
}

func (a *AuthServiceImpl) Validate(jwtToken string) (*types.Claims, bool) {
	customClaims := &types.Claims{}

	token, err := jwt.ParseWithClaims(jwtToken, customClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.signingKey), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}

	return customClaims, true
}
