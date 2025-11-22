package user_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
	"fmt"
)

func (r *userRepository) UserExists(ctx context.Context, q repositories.QueryExecer, userID string) (bool, error) {

	var exists bool
	err := q.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM users WHERE ID = $1
		)`,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check user's %s existence: %w", userID, err)
	}

	return exists, nil

}
