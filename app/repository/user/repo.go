package user

import (
	"broker/app/models"
	"broker/app/service"
	"broker/pkg/pg"
	"database/sql"
	"errors"
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

func (r *Repository) GetOne(email string) (*models.User, error) {
	var user models.User

	row := sq.Select("id","password", "email", "first_name", "last_name").
		From("users").
		Where(sq.Eq{"email": email}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email, &user.FirstName, &user.LastName); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, service.UserNotFoundErr
		}
		return nil, err
	}

	return &user, nil
}
