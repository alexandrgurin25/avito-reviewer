package prserv

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	prrepo "avito-reviewer/internal/repositories/pull_request_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"

	"context"
)

type PRService interface {
	CreatePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error)
	MergePR(ctx context.Context, id string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, pr *models.ReasignPR) (*models.PullRequest, string, error)
	GetReviewerStats(ctx context.Context) (map[string]int, error)
}

type prService struct {
	userRepo userrepo.Repository
	db       repositories.DB
	prRepo   prrepo.Repository
}

func NewService(userRepo userrepo.Repository,
	prRepo prrepo.Repository, db repositories.DB) PRService {
	return &prService{userRepo: userRepo, prRepo: prRepo, db: db}
}
