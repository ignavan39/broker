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

func (r *Repository) Create(nickname string, email string, password string, lastName string, firstName string) (*models.User, error) {
	user := &models.User{}

	row := sq.Insert("users").
		Columns("password", "email", "nickname", "first_name", "last_name").
		Values(password, email, nickname, firstName, lastName).
		Suffix("returning id, password, email, nickname, first_name,last_name, avatar_url").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email, &user.Nickname, &user.FirstName, &user.LastName, &user.AvatarURL); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			return nil, service.DuplicateUserErr
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetOneByEmail(email string) (*models.User, error) {
	var user models.User

	row := sq.Select("id", "password", "email", "first_name", "last_name", "avatar_url", "nickname").
		From("users").
		Where(sq.Eq{"email": email}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.AvatarURL, &user.Nickname); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, service.UserNotFoundErr
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetOneByNickname(nickname string) (*models.User, error) {
	var user models.User

	row := sq.Select("id", "password", "email", "first_name", "last_name", "avatar_url", "nickname").
		From("users").
		Where(sq.Eq{"nickname": nickname}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&user.Id, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.AvatarURL, &user.Nickname); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, service.UserNotFoundErr
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetEmailById(userId string) (string, error) {
	var email string

	row := sq.Select("email").
		From("users").
		Where(sq.Eq{"id": userId}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", service.UserNotFoundErr
		}
		return "", err
	}

	return email, nil
}
