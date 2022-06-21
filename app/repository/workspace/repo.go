package workspace

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

func (r *Repository) Create(email string, name string, isPrivate bool) (*models.Workspace, error) {
	var workspace models.Workspace

	row := sq.Insert("workspaces").
		Columns("name", "is_private").
		Values(name, isPrivate).
		Suffix("returning id, name, created_at, is_private").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			return nil, service.DuplicateWorkspaceErr
		}
		return nil, err
	}

	_, err := sq.Insert("workspace_accesses").
		Columns("workspace_id", "email", `"type"`).
		Values(workspace.ID, email, models.ADMIN).
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			return nil, service.DuplicateWorkspaceEmailErr
		}
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
			return service.WorkspaceNotExistsErr
		}
		return err
	}
	return nil
}

func (r *Repository) Update(workspaceID string, name *string, isPrivate *bool) (*models.Workspace, error) {
	var workspace models.Workspace

	qb := sq.Update("workspaces")

	if name != nil {
		qb =qb.Set("name", *name)
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
			return nil, service.WorkspaceNotExistsErr
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
		InnerJoin("users u ON wa.email = u.email").
		Where(sq.Eq{"u.id": id}).
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

func (r *Repository) GetAccessByUserId(userID string, workspaceID string) (*models.Workspace, error) {
	var workspace models.Workspace

	row := sq.Select("w.id", "w.name", "w.created_at", "w.is_private").
		From("workspaces w").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = w.id AND wa.type = ?", models.ADMIN).
		InnerJoin("users u ON u.email = wa.email").
		Where(sq.Eq{"w.id": workspaceID, "u.id": userID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.WorkspaceNotExistsErr
		}

		return nil, err
	}

	return &workspace, nil
}
