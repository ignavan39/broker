package invitation

import (
	"broker/core/models"
	"broker/core/service"
	"broker/pkg/pg"
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	
	blogger "github.com/sirupsen/logrus"
)

type Repository struct {
	pool pg.Pool
}

func NewRepository(pool pg.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) CheckInvites(userID string, email string) error {
	blogger.Println("HERE")

	tx, err := r.pool.Write().Begin()

	if err != nil {
		return err
	}

	rows, err := sq.Update("invitations").
		Set("status", models.ACCEPTED).
		Where(sq.Eq{"ricipient_email": email}).
		Suffix("returning workspace_id").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		Query()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = tx.Commit(); err != nil {
				return err
			}

			return nil
		}

		if err = tx.Rollback(); err != nil {
			return err
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
			RunWith(tx).
			PlaceholderFormat(sq.Dollar).
			Exec()
		if err != nil {
			duplicate := strings.Contains(err.Error(), "duplicate")

			if duplicate {
				if err = tx.Commit(); err != nil {
					return err
				}

				return service.DuplicateUserErr
			}

			if err = tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) SendInvitation(senderID string, workspaceID string, ricipientEmail string) (*models.Invitation, error) {
	var invitation models.Invitation

	row := sq.Select("id", "email", "nickname", "first_name", "last_name", "avatar_url").
		From("users").
		Where(sq.Eq{"id": senderID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.Sender.ID, &invitation.Sender.Email, &invitation.Sender.Nickname, 
						&invitation.Sender.FirstName, &invitation.Sender.LastName, &invitation.Sender.AvatarURL); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, service.UserNotFoundErr
		}
		return nil, err
	}

	row = sq.Insert("invitations").
		Columns("sender_id", "ricipient_email", "workspace_id").
		Values(senderID, ricipientEmail, workspaceID).
		Suffix("returning id, ricipient_email, workspace_id, status").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.ID, &invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status); err != nil {
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

	rows, err := sq.Select("i.id", "i.sender_id", "u.email", "u.nickname",
		"u.first_name", "u.last_name", "u.avatar_url", "i.ricipient_email",
		"i.workspace_id", "i.status").
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

		if err := rows.Scan(&invitation.ID, &invitation.Sender.ID,
			&invitation.Sender.Email, &invitation.Sender.Nickname,
			&invitation.Sender.FirstName, &invitation.Sender.LastName, &invitation.Sender.AvatarURL,
			&invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status); err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (r *Repository) CancelInvitation(invitationID string) (*models.Invitation, error) {
	var invitation models.Invitation

	row := sq.Update("invitations").
		Set(`"status"`, models.CANCELED).
		Where(sq.Eq{"id": invitationID}).
		Suffix("returning id, sender_id, ricipient_email, workspace_id, status").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.ID, &invitation.Sender.ID, &invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.InvitationNotFoundErr
		}

		return nil, err
	}

	row = sq.Select("email", "nickname", "first_name", "last_name", "avatar_url").
		From("users").
		Where(sq.Eq{"id": invitation.Sender.ID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.Sender.Email, &invitation.Sender.Nickname, 
						&invitation.Sender.FirstName, &invitation.Sender.LastName, &invitation.Sender.AvatarURL); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, service.UserNotFoundErr
		}
		return nil, err
	}

	return &invitation, nil
}
