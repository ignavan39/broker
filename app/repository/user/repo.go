package user

import (
	"broker/app/models"
	"broker/app/service"
	"broker/pkg/pg"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	pool pg.Pool
}

func NewRepository(pool pg.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(email string, password string, lastName string, firstName string) (*models.User, error) {
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
			return nil, service.DuplicateUserErr
		}
		return nil, err
	}

	return user, nil
}
