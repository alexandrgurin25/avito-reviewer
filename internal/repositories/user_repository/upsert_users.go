package userrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"

	"github.com/jackc/pgx/v5"
)

func (*userRepository) UpsertUsers(ctx context.Context, q repositories.QueryExecer, t *models.Team) (*models.Team, error) {
	var res models.Team

	batch := &pgx.Batch{}

	for _, u := range t.Members {
		batch.Queue(`
        INSERT INTO users (id, username, is_active, team_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            is_active = EXCLUDED.is_active,
            team_id = EXCLUDED.team_id
    `, u.ID, u.Username, u.IsActive, t.ID)
		res.Members = append(res.Members, u)
	}

	br := q.SendBatch(ctx, batch)
	err := br.Close()
	if err != nil {
		return nil, err
	}

	return &res, nil
}
