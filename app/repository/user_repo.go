package repository

import (
	"broker/app/models"
)

type UserRepository interface {
	Create(email string, password string, lastName string, firstName string) (*models.User, error)
}
