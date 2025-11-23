package pull_request_services

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories/pull_request_repository"
	"avito-reviewer/internal/repositories/user_repository"
	"context"
)

type PRService interface {
	CreatePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error)
	MergePR(ctx context.Context, id string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, pr *models.ReasignPR) (*models.PullRequest, string, error)
}

type prService struct {
	userRepo user_repository.Repository

	prRepo pull_request_repository.Repository
}

func NewService(userRepo user_repository.Repository,
	prRepo pull_request_repository.Repository) PRService {
	return &prService{userRepo: userRepo, prRepo: prRepo}
}
