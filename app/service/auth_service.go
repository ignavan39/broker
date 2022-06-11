package service

import (
	"broker/app/types"
)

type AuthService interface {
	Refresh(id string) (map[string]string, error)
	Validate(jwtToken string) (*types.Claims, bool)
}
