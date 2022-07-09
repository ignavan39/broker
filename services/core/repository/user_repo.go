package repository

import (
	"broker/core/models"
)

type UserRepository interface {
	Create(nickname string, email string, password string, lastName string, firstName string) (*models.User, error)
	GetOneByEmail(email string) (*models.User, error)
	GetOneByNickname(nickname string) (*models.User, error)
	GetEmailById(userID string) (string, error)
	GetOneById(id string) (*models.User, error)
}
