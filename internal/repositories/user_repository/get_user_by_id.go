package user_repository

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *userRepository) GetByID(ctx context.Context, q repositories.QueryExecer, id string) (*models.User, error) {
	var u models.User
	err := q.QueryRow(ctx,
		`SELECT id, username, team_id, is_active
         FROM users WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
