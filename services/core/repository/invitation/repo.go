package invitation

import (
	"broker/core/config"
	"broker/core/models"
	"broker/core/service"
	"broker/pkg/pg"
	"broker/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
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

func (r *Repository) AcceptInvitation(userID string, code string) error {
	tx, err := r.pool.Write().Begin()
	var workspaceID string

	if err != nil {
		return err
	}

	row := sq.Update("invitations").
		Set("status", models.ACCEPTED).
		Where(sq.Eq{"code": code}).
		Suffix("returning workspace_id").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&workspaceID); err != nil {
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

	_, err = sq.Insert("workspace_accesses").
		Columns("workspace_id", "user_id").
		Values(workspaceID, userID).
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

	code := utils.CryptString(fmt.Sprintf("%s%s", senderID, workspaceID), config.GetConfig().Invitation.InvitationHashSalt)

	row = sq.Insert("invitations").
		Columns("sender_id", "ricipient_email", "workspace_id, code, system_status").
		Values(senderID, ricipientEmail, workspaceID, code, models.SEND).
		Suffix("returning id, created_at, ricipient_email, workspace_id, status, system_status, code").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.ID, &invitation.CreatedAt, &invitation.RicipientEmail,
		&invitation.WorkspaceID, &invitation.Status, &invitation.SystemStatus,
		&invitation.Code); err != nil {
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

	rows, err := sq.Select("i.id", "i.created_at", "i.sender_id", "u.email", "u.nickname",
		"u.first_name", "u.last_name", "u.avatar_url", "i.ricipient_email",
		"i.workspace_id", "i.status").
		From("invitations i").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = i.workspace_id").
		InnerJoin("users u ON wa.user_id = u.id").
		Where(sq.Eq{"u.id": userID, "i.workspace_id": workspaceID}).
		OrderBy("created_at DESC", "sender_id DESC").
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

		if err := rows.Scan(&invitation.ID, &invitation.CreatedAt, &invitation.Sender.ID,
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

	_, err := sq.Select("u.id").
		From("users u").
		InnerJoin("workspace_accesses wa ON wa.user_id = u.id").
		Where(sq.NotEq{"wa.`type`": models.USER}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		Query()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.WorkspaceAccessDeniedErr
		}

		return nil, err
	}

	row := sq.Update("invitations").
		Set(`"status"`, models.CANCELED).
		Where(sq.Eq{"id": invitationID}).
		Suffix("returning id, created_at, sender_id, ricipient_email, workspace_id, status, system_status, code").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&invitation.ID, &invitation.CreatedAt, &invitation.Sender.ID,
		&invitation.RicipientEmail, &invitation.WorkspaceID, &invitation.Status,
		&invitation.SystemStatus, &invitation.Code); err != nil {
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
