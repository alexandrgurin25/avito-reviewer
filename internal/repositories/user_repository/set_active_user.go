package user_repository

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *userRepository) SetActive(ctx context.Context, q repositories.QueryExecer, id string, active bool) (*models.User, int, error) {

	var user models.User

	var teamID int

	err := q.QueryRow(
		ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, username, team_id, is_active `,
		active, id,
	).Scan(&user.ID, &user.Username, &teamID, &user.IsActive)

	if err != nil {
		return nil, 0, err
	}

	return &user, teamID, nil
}
