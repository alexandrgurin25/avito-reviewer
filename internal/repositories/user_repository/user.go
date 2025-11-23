package user_repository

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

type Repository interface {
	UpsertUsers(ctx context.Context, q repositories.QueryExecer, t *models.Team) (*models.Team, error)
	GetUsersByTeam(ctx context.Context, db repositories.QueryExecer, teamID int) (*models.Team, error)
	GetExistingUsers(ctx context.Context, db repositories.QueryExecer, ids []string) (map[string]string, error)
	SetActive(ctx context.Context, q repositories.QueryExecer, id string, active bool) (*models.User, int, error)
	UserExists(ctx context.Context, q repositories.QueryExecer, userID string) (bool, error)
	GetByID(ctx context.Context, q repositories.QueryExecer, id string) (*models.User, error)

	GetRandomReviewers(ctx context.Context, q repositories.QueryExecer, team, exclude string) ([]string, error)
}

type userRepository struct {
	pool repositories.DB
}

func NewUserRepository(pool repositories.DB) Repository {
	return &userRepository{pool: pool}
}
