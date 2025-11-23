package prrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"time"
)

type Repository interface {
	BeginTx(ctx context.Context) (repositories.Tx, error)
	Create(ctx context.Context, q repositories.QueryExecer, pr *models.PullRequest) error
	Exists(ctx context.Context, q repositories.QueryExecer, id string) (bool, error)
	GetByID(ctx context.Context, q repositories.QueryExecer, id string) (*models.PullRequest, error)
	SetStatusOnMerged(ctx context.Context, q repositories.QueryExecer, id string, mergedAt time.Time) error
	GetIDByReviewerID(ctx context.Context, q repositories.QueryExecer, userID string) ([]string, error)
	GetPRsByID(ctx context.Context, db repositories.QueryExecer, ids []string) ([]models.PullRequest, error)
	ReplaceReviewer(ctx context.Context, q repositories.QueryExecer, prID, oldID, newID string) error
}

type prRepository struct {
	pool repositories.DB
}

func NewPRRepository(pool repositories.DB) Repository {
	return &prRepository{pool: pool}
}

func (r *prRepository) BeginTx(ctx context.Context) (repositories.Tx, error) {
	return r.pool.BeginTx(ctx)
}
