package user_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (u *userRepository) GetExistingUsers(ctx context.Context,
	db repositories.QueryExecer, ids []string) (map[string]string, error) {

	rows, err := db.Query(ctx,
		`SELECT id, team_id FROM users WHERE id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string]string)
	for rows.Next() {
		var id, team string
		if err := rows.Scan(&id, &team); err != nil {
			return nil, err
		}
		res[id] = team
	}
	return res, nil
}
