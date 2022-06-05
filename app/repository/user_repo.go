package repository

import (
	"broker/app/models"
	"broker/pkg/pg"
	sq "github.com/Masterminds/squirrel"
)

type UserRepository interface{
	Create(email string, password string) (*models.User, error)
} 

type UserRepositoryImpl struct {
	pool pg.Pool
}

func NewUserRepository(pool pg.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		pool: pool,
	}
}

func (r *UserRepositoryImpl) Create(email string, password string) (*models.User, error) {
	user := &models.User{}

	row := sq.Insert("users").
		Columns("password", "email").
		Values(password, email).
		Suffix("returning id, password, email").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email); err != nil {
		return nil, err
	}

	return user, nil
}