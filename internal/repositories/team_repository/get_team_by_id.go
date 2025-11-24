package teamrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"fmt"
)

func (*teamRepository) GetTeamNameByID(ctx context.Context, q repositories.QueryExecer, id int) (*models.Team, error) {

	var res models.Team
	err := q.QueryRow(
		ctx,
		`SELECT name FROM teams WHERE id = $1`,
		id,
	).Scan(&res.Name)

	if err != nil {
		return nil, fmt.Errorf("failed to find team's %d Name: %w", id, err)
	}

	return &res, nil

}
