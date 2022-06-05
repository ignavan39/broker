package services

import "context"

type AuthService interface {
	CreateToken(ctx context.Context, id string) (string, error)
}

type AuthServiceImpl struct {
	
}