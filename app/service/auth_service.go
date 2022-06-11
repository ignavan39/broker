package service

import (
	"broker/app/dto"

	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	Id string `json:"id"`
}

type AuthService interface {
	SignUp(payload dto.SignUpPayload) (*dto.SignResponse, error)
	Validate(jwtToken string) (*Claims, bool)
}
