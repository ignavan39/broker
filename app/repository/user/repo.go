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
	if err := row.Scan(&user.ID, &user.Password, &user.Email, &user.Nickname, &user.FirstName, &user.LastName, &user.AvatarURL); err != nil {
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
	if err := row.Scan(&user.ID, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.AvatarURL, &user.Nickname); err != nil {
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
	if err := row.Scan(&user.ID, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.AvatarURL, &user.Nickname); err != nil {
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

func (r *Repository) CheckInvites(userID string, email string) error {
	rows, err := sq.Update("invitations").
		Set("status", models.ACCEPTED).
		Where(sq.Eq{"ricipient_email": email}).
		Suffix("returning workspace_id").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		Query()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	defer rows.Close()

	for rows.Next() {
		var workspace_id string

		if err := rows.Scan(&workspace_id); err != nil {
			return err
		}

		_, err := sq.Insert("workspace_accesses").
			Columns("workspace_id", "user_id").
			Values(workspace_id, userID).
			RunWith(r.pool.Write()).
			PlaceholderFormat(sq.Dollar).
			Exec()
		if err != nil {
			duplicate := strings.Contains(err.Error(), "duplicate")

			if duplicate {
				return service.DuplicateWorkspaceAccessErr
			}

			return err
		}
	}

	return nil
}

func (r *Repository) SendInvitation(senderID string, workspaceID string, ricipientEmail string) (*models.Invitation, error) {
	var invitation models.Invitation

	row := sq.Insert("invitations").
		Columns("sender_id", "ricipient_email", "workspace_id").
		Values(senderID, ricipientEmail, workspaceID).
		Suffix("returning id, sender_id, ricipient_email, workspace_id, status").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.ID, &invitation.SenderID, &invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")

		if duplicate {
			return nil, service.DuplicateInvitationErr
		}

		return nil, err
	}

	return &invitation, nil
}

func (r *Repository) GetInvitationsByWorkspaceID(userID string, workspaceID string) ([]models.Invitation, error) {
	invitations := make([]models.Invitation, 0)

	rows, err := sq.Select("i.id", "i.sender_id", "i.ricipient_email", "i.workspace_id", "i.status").
		From("invitations i").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = i.workspace_id").
		InnerJoin("users u ON wa.user_id = u.id").
		Where(sq.Eq{"u.id": userID, "i.workspace_id": workspaceID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return invitations, nil
		}
		return invitations, err
	}

	defer rows.Close()

	for rows.Next() {
		var invitation models.Invitation

		if err := rows.Scan(&invitation.ID, &invitation.SenderID, &invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status); err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}
