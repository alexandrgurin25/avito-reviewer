package team_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
	"fmt"
)

func (r *teamRepository) TeamExists(ctx context.Context, q repositories.QueryExecer, name string) (bool, error) {

	var exists bool
	err := q.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM teams WHERE name = $1
		)`,
		name,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check team's %s existence: %w", name, err)
	}

	return exists, nil

}
