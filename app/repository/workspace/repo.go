package workspace

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

func (r *Repository) Create(email string, name string, isPrivate bool) (*models.Workspace, error) {
	var workspace models.Workspace

	row := sq.Insert("workspaces").
		Columns("name", "is_private").
		Values(name, isPrivate).
		Suffix("returning id, name, created_at, is_private").
		RunWith(r.pool.Write()).
		PlaceholderFormat(sq.Dollar).
		QueryRow()
	if err := row.Scan(&workspace.Id, &workspace.Name, &workspace.CreatedAt, &workspace.IsPrivate); err != nil {
		duplicate := strings.Contains(err.Error(), "duplicate")
		if duplicate {
			return nil, service.DuplicateWorkspaceErr
		}
		return nil, err
	}

	_, err := sq.Insert("workspace_accesses").
		Columns("workspace_id", "email", `"type"`).
		Values(workspace.Id, email, models.ADMIN).
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
