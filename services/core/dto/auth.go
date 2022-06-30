package dto

import (
	"broker/core/models"
	"errors"
	"regexp"
	"time"
)

type SignPayloadBase struct {
	Password string  `json:"password"`
	Email    *string `json:"email,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
}

type SignUpPayload struct {
	SignPayloadBase
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Code      int    `json:"code"`
}

type TokenResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expireAt"`
}
type AuthResponse struct {
	Access  TokenResponse `json:"access"`
	Refresh TokenResponse `json:"refresh"`
}

type SignResponse struct {
	User models.User  `json:"user"`
	Auth AuthResponse `json:"auth"`
}

type SignPayloadResponseBuilder struct {
	payload SignResponse
}

func NewSignPayloadResponseBuilder() *SignPayloadResponseBuilder {
	return &SignPayloadResponseBuilder{
		payload: SignResponse{
			User: models.User{},
			Auth: AuthResponse{},
		},
	}
}

func (sprb *SignPayloadResponseBuilder) WithUser(user models.User) *SignPayloadResponseBuilder {
	sprb.payload.User = user
	return sprb
}

func (sprb *SignPayloadResponseBuilder) WithAccessToken(accessToken string, duration time.Duration) *SignPayloadResponseBuilder {
	sprb.payload.Auth.Access = TokenResponse{
		Token:    accessToken,
		ExpireAt: time.Now().Add(duration),
	}
	return sprb
}

func (sprb *SignPayloadResponseBuilder) WithRefreshToken(refreshToken string, duration time.Duration) *SignPayloadResponseBuilder {
	sprb.payload.Auth.Refresh = TokenResponse{
		Token:    refreshToken,
		ExpireAt: time.Now().Add(duration),
	}
	return sprb
}

func (sprb *SignPayloadResponseBuilder) Exec() SignResponse {
	return sprb.payload
}

type SendCodePayload struct {
	Email string `json:"email"`
}

func (scp *SendCodePayload) Validate() error {
	if !isCorrectEmail(scp.Email) {
		return errors.New("email must be not empty string")
	}
	return nil
}

type SignInPayload = SignPayloadBase

func isCorrectEmail(email string) bool {
	pattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	matched, _ := regexp.Match(pattern, []byte(email))
	return matched
}

func (p *SignUpPayload) Validate() error {
	if len(p.Password) < 6 {
		return errors.New("password too short")
	}

	if p.Email == nil || !isCorrectEmail(*p.Email) {
		return errors.New("email must be not empty string")
	}

	if p.Nickname == nil || len(*p.Nickname) == 0 {
		return errors.New("nickname must be not empty string")
	}

	if len(p.LastName) == 0 {
		return errors.New("last name must be not empty string")
	}

	if len(p.FirstName) == 0 {
		return errors.New("first name must be not empty string")
	}

	return nil
}

func (p *SignInPayload) Validate() error {
	if len(p.Password) < 6 {
		return errors.New("password too short")
	}

	if p.Email == nil || !isCorrectEmail(*p.Email) {
		return errors.New("email must be not empty string")
	}

	if p.Nickname == nil || len(*p.Nickname) == 0 {
		return errors.New("nickname must be not empty string")
	}

	return nil
}
