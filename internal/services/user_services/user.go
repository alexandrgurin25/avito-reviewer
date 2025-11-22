package user_services

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"avito-reviewer/internal/repositories/team_repository"
	"avito-reviewer/internal/repositories/user_repository"
	"context"
)

type UserService interface {
	SetIsActiveUser(ctx context.Context, u *models.User) (*models.User, error)
}

type userService struct {
	userRepo user_repository.Repository
	teamRepo team_repository.Repository
	db       repositories.DB
}

func NewService(userRepo user_repository.Repository, teamRepo team_repository.Repository, db repositories.DB) UserService {
	return &userService{userRepo: userRepo, teamRepo: teamRepo, db: db}
}
