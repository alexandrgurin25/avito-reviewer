package userrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

func (*userRepository) GetUsersByTeam(ctx context.Context,
	db repositories.QueryExecer, teamID int) (*models.Team, error) {

	rows, err := db.Query(ctx,
		`SELECT id, username, is_active FROM users WHERE team_id = $1`,
		teamID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := models.Team{}
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.IsActive); err != nil {
			return nil, err
		}

		res.Members = append(res.Members, user)
	}

	return &res, nil
}
