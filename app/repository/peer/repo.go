package peer

import (
	"broker/app/models"
	"broker/pkg/pg"
	"database/sql"
	"errors"

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

func (r *Repository) GetMany(userID string, workspaceID string) ([]models.Peer, error) {
	peers := make([]models.Peer, 0)

	rows, err := sq.Select("p.id", "p.name", "p.created_at", "wa.workspace_id").
		From("peers p").
		InnerJoin("workspace_accesses wa ON wa.workspace_id = p.workspace_id").
		InnerJoin("users u ON u.id = wa.user_id").
		Where(sq.Eq{"wa.workspace_id": workspaceID, "u.id": userID}).
		RunWith(r.pool.Read()).
		PlaceholderFormat(sq.Dollar).
		Query()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return peers, nil
		}
		return peers, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Peer

		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.WorkspaceId); err != nil {
			return nil, err
		}

		peers = append(peers, p)
	}
	return peers, nil
}
