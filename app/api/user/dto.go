package user

import (
	"broker/app/models"
	"errors"
)

type SignUpPayload struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Code      int    `json:"code"`
}

func (p *SignUpPayload) Validate() error {
	if (len(p.Password) == 0){
		return errors.New("password too short")
	}
	if (len(p.Email) == 0) {
		return errors.New("email must be not empty string")
	}
	if (len(p.LastName) == 0) {
		return errors.New("last name must be not empty string")
	}
	if (len(p.FirstName) == 0) {
		return errors.New("first name must be not empty string")
	}
	return nil
}

type SignResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

type SendCodePayload struct {
	Email string `json:"email"`
}
