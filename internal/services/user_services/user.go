package userserv

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	prrepo "avito-reviewer/internal/repositories/pull_request_repository"
	teamRepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"

	"context"
)

type UserService interface {
	SetIsActiveUser(ctx context.Context, u *models.User) (*models.User, error)
	GetReview(ctx context.Context, userID string) (*models.UserPR, error)
}

type userService struct {
	userRepo userrepo.Repository
	teamRepo teamRepo.Repository
	prRepo   prrepo.Repository
	db       repositories.DB
}

func NewService(userRepo userrepo.Repository, teamRepo teamRepo.Repository, prRepo prrepo.Repository, db repositories.DB) UserService {
	return &userService{userRepo: userRepo, teamRepo: teamRepo, prRepo: prRepo, db: db}
}
