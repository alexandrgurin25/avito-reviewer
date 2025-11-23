package teamrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"fmt"
)

func (*teamRepository) CreateTeam(ctx context.Context, q repositories.QueryExecer, name string) (*models.Team, error) {

	res := models.Team{}

	err := q.QueryRow(
		ctx,
		`INSERT INTO teams (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&res.ID)

	res.Name = name

	if err != nil {
		return nil, fmt.Errorf("falied team's %s create error %v", res.Name, err)
	}

	return &res, nil
}
