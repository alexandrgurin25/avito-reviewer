package user_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *userRepository) GetReassignCandidates(
	ctx context.Context,
	q repositories.QueryExecer,
	team string,
	assigned []string,
	oldUser string,
	author string,
) ([]string, error) {

	rows, err := q.Query(ctx, 
		`SELECT id
        FROM users
        WHERE team_id = $1
          AND is_active = TRUE
          AND id != $2          
          AND id != $3         
          AND id <> ALL($4)     
        ORDER BY random()`,
		team, oldUser, author, assigned,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}
