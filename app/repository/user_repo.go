package repository

import (
	"broker/app/models"
	"broker/pkg/pg"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository interface {
	Create(email string, password string, lastName string, firstName string) (*models.User, error)
}

type UserRepositoryImpl struct {
	pool pg.Pool
}

func NewUserRepository(pool pg.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		pool: pool,
	}
}

func (r *UserRepositoryImpl) Create(email string, password string, lastName string, firstName string) (*models.User, error) {
	user := &models.User{}

	row := sq.Insert("users").
		Columns("password", "email", "first_name", "last_name").
		Values(password, email, firstName, lastName).
		Suffix("returning id, password, email, first_name,last_name").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email, &user.FirstName, &user.LastName); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			return nil, DuplicateUserErr
		}
		return nil, err
	}

	return user, nil
}
