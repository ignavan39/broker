package repository

import (
	"broker/app/models"
)

type UserRepository interface {
	Create(nickname string, email string, password string, lastName string, firstName string) (*models.User, error)
	GetOneByEmail(email string) (*models.User, error)
	GetOneByNickname(nickname string) (*models.User, error)
}
