package dto

import (
	"broker/app/models"
	"errors"
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

func (p *SignUpPayload) Validate() error {
	if len(p.Password) == 0 {
		return errors.New("password too short")
	}
	if p.Email == nil || len(*p.Email) == 0 {
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

type SignInPayload = SignPayloadBase

func (p *SignInPayload) Validate() error {
	if len(p.Password) == 0 {
		return errors.New("password too short")
	}
	if (p.Email == nil || len(*p.Email) == 0) && (p.Nickname == nil || len(*p.Nickname) == 0) {
		return errors.New("email/nickname must be not empty string")
	}

	return nil
}

type SignResponse struct {
	User models.User       `json:"user"`
	Auth map[string]string `json:"auth"`
}

type SendCodePayload struct {
	Email string `json:"email"`
}