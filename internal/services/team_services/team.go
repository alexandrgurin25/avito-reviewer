package teamserv

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"

	teamRepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"
	"context"
)

type TeamService interface {
	AddTeam(ctx context.Context, t *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type teamService struct {
	userRepo userrepo.Repository
	teamRepo teamRepo.Repository
	db       repositories.DB
}

func NewService(userRepo userrepo.Repository,
	teamRepo teamRepo.Repository, db repositories.DB) TeamService {
	return &teamService{userRepo: userRepo, teamRepo: teamRepo, db: db}
}
