package user_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *userRepository) GetRandomReviewers(ctx context.Context, q repositories.QueryExecer, teamID, exclude string) ([]string, error) {
	rows, err := q.Query(ctx,
		`SELECT id
         FROM users
         WHERE team_id = $1 AND is_active = TRUE AND id != $2
         ORDER BY random()
         LIMIT 2`,
		teamID, exclude,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	return res, nil
}
