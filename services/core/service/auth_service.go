package service

import (
	"broker/core/dto"
	"context"

	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	Id string `json:"id"`
}

type AuthService interface {
	SignUp(ctx context.Context, payload dto.SignUpPayload) (*dto.SignResponse, error)
	SignIn(payload dto.SignInPayload) (*dto.SignResponse, error)
	SendVerifyCode(ctx context.Context, email string) error
	Validate(jwtToken string) (*Claims, bool)
	Refresh(jwtToken string) (*dto.SignResponse, error)
}
