package teamrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

type Repository interface {
	BeginTx(ctx context.Context) (repositories.Tx, error)

	GetTeamIDByName(ctx context.Context, q repositories.QueryExecer, name string) (*models.Team, error)
	GetTeamNameByID(ctx context.Context, q repositories.QueryExecer, id int) (*models.Team, error)

	TeamExists(ctx context.Context, q repositories.QueryExecer, name string) (bool, error)
	CreateTeam(ctx context.Context, q repositories.QueryExecer, name string) (*models.Team, error)
}

type teamRepository struct {
	pool repositories.DB
}

func NewTeamRepository(pool repositories.DB) Repository {
	return &teamRepository{pool: pool}
}

func (r *teamRepository) BeginTx(ctx context.Context) (repositories.Tx, error) {
	return r.pool.BeginTx(ctx)
}
