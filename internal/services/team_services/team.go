package team_services

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"avito-reviewer/internal/repositories/team_repository"
	"avito-reviewer/internal/repositories/user_repository"
	"context"
)

type TeamService interface {
	AddTeam(ctx context.Context, t *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type teamService struct {
	userRepo user_repository.Repository
	teamRepo team_repository.Repository
	db       repositories.DB
}

func NewService(userRepo user_repository.Repository,
	teamRepo team_repository.Repository, db repositories.DB) TeamService {
	return &teamService{userRepo: userRepo, teamRepo: teamRepo, db: db}
}
