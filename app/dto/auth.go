package dto

import (
	"broker/app/models"
	"errors"
	"regexp"
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

type SignResponse struct {
	User models.User       `json:"user"`
	Auth map[string]string `json:"auth"`
}

type SendVerifyCodePayload struct {
	Email string `json:"email"`
}

type VerifyCodePayload struct {
	Code string `json:"code"`
}

type SendCodePayload struct {
	Email string `json:"email"`
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

	if p.Email == nil || isCorrectEmail(*p.Email) {
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

	if p.Email == nil || isCorrectEmail(*p.Email) {
		return errors.New("email must be not empty string")
	}

	if p.Nickname == nil || len(*p.Nickname) == 0 {
		return errors.New("nickname must be not empty string")
	}

	return nil
}

func (p *SendVerifyCodePayload) Validate() error {
	if isCorrectEmail(p.Email) {
		return errors.New("email is not correct")
	}

	return nil
}

func (p *VerifyCodePayload) Validate() error {
	if len(p.Code) != 5 {
		return errors.New("code is not correct")
	}

	return nil
}
