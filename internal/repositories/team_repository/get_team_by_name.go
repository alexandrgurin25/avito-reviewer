package teamrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"fmt"
)

func (*teamRepository) GetTeamIDByName(ctx context.Context, q repositories.QueryExecer, name string) (*models.Team, error) {

	var res models.Team
	err := q.QueryRow(
		ctx,
		`SELECT id FROM teams WHERE name = $1`,
		name,
	).Scan(&res.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to find team's %s ID: %w", name, err)
	}

	return &res, nil

}
