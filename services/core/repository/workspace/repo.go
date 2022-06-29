package workspace

import (
	"broker/core/models"
	"broker/core/service"
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

func (r *Repository) Create(userID string, name string, isPrivate bool) (*models.Workspace, error) {
	var workspace models.Workspace

	tx, err := r.pool.Write().Begin()

	if err != nil {
		return nil, err
	}

	row := sq.Insert("workspaces").
		Columns("name", "is_private").
		Values(name, isPrivate).
		Suffix("returning id, name, created_at, is_private").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")

		if duplicate {
			if err = tx.Commit(); err != nil {
				return nil, err
			}

			return nil, service.DuplicateWorkspaceErr
		}

		if err = tx.Rollback(); err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = sq.Insert("workspace_accesses").
		Columns("workspace_id", "user_id", `"type"`).
		Values(workspace.ID, userID, models.ADMIN).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			if err = tx.Commit(); err != nil {
				return nil, err
			}

			return nil, service.DuplicateWorkspaceEmailErr
		}

		if err = tx.Rollback(); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (r *Repository) Delete(workspaceID string) error {
	_, err := sq.Delete("workspaces").
		Where(sq.Eq{"id": workspaceID}).
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.WorkspaceAccessDeniedErr
		}
		return err
	}
	return nil
}

func (r *Repository) Update(workspaceID string, name *string, isPrivate *bool) (*models.Workspace, error) {
	var workspace models.Workspace

	qb := sq.Update("workspaces")

	if name != nil {
		qb = qb.Set("name", *name)
	}

	if isPrivate != nil {
		qb = qb.Set("is_private", *isPrivate)
	}
	qb = qb.Where(sq.Eq{"id": workspaceID}).
		Suffix("returning id, name, created_at, is_private").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar)
	row := qb.QueryRow()

	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.WorkspaceAccessDeniedErr
		}

		return nil, err
	}

	return &workspace, nil
}

func (r *Repository) GetManyByUserId(id string) ([]models.Workspace, error) {
	workspaces := make([]models.Workspace, 0)

	rows, err := sq.Select("w.id", `w."name"`, "w.created_at", "w.is_private").
		From("workspaces w").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = w.id").
		InnerJoin("users u ON wa.user_id = u.id").
		Where(sq.Eq{"u.id": id}).
		OrderBy("w.is_private DESC", "w.created_at DESC").
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workspaces, nil
		}
		return workspaces, err
	}

	defer rows.Close()

	for rows.Next() {
		var w models.Workspace

		if err := rows.Scan(&w.ID, &w.Name, &w.CreatedAt, &w.IsPrivate); err != nil {
			return nil, err
		}

		workspaces = append(workspaces, w)
	}
	return workspaces, nil
}

func (r *Repository) GetWorkspaceByUserId(userID string, workspaceID string) (*models.Workspace, error) {
	var workspace models.Workspace

	row := sq.Select("w.id", "w.name", "w.created_at", "w.is_private").
		From("workspaces w").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = w.id").
		InnerJoin("users u ON u.id = wa.user_id").
		Where(sq.Eq{"w.id": workspaceID, "u.id": userID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.WorkspaceAccessDeniedErr
		}

		return nil, err
	}

	return &workspace, nil
}

func (r *Repository) GetAccessByUserId(userID string, workspaceID string) (*string, error) {
	var accessType string

	row := sq.Select("wa.type").
		From("workspace_accesses wa").
		InnerJoin("users u ON u.id = wa.user_id").
		Where(sq.Eq{"u.id": userID, "wa.workspace_id": workspaceID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&accessType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.WorkspaceAccessDeniedErr
		}

		return nil, err
	}

	return &accessType, nil
}

func (r *Repository) GetWorkspaceUsersCount(workspaceID string) (int, error) {
	var usersCount int

	row := sq.Select("COUNT(*)").
		From("workspace_accesses wa").
		Where(sq.Eq{"wa.workspace_id": workspaceID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&usersCount); err != nil {
		return 0, err
	}

	if usersCount == 0 {
		return usersCount, service.WorkspaceAccessDeniedErr
	}

	return usersCount, nil
}

func (r *Repository) ChangeUserRole(role string, userID string, workspaceID string) error {
	_, err := sq.Update("workspace_accesses").
		Set("role", role).
		Where(sq.Eq{"user_id": userID, "workspace_id": workspaceID}).
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.UserNotFoundErr
		}

		return err
	}

	return nil
}

func (r *Repository) BanUser(bannedUserID string, workspaceID string) error {
	_, err := sq.Delete("workspace_accesses").
		Where(sq.Eq{"user_id": bannedUserID, "workspace_id": workspaceID}, sq.NotEq{`"type"`: models.ADMIN}).
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.UserNotFoundErr
		}

		return err
	}

	return nil
}
