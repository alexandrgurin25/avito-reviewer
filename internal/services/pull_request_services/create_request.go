package prserv

import (
	"avito-reviewer/internal/models"
	"context"
	"time"
)

func (s *prService) CreatePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {

	tx, err := s.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверить существует ли такой пользователь
	hasUser, err := s.userRepo.UserExists(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	if !hasUser {
		return nil, models.ErrNotFound
	}

	// Проверить существует ли pull_request
	hasPR, err := s.prRepo.Exists(ctx, tx, pr.ID)
	if err != nil {
		return nil, err
	}
	if hasPR {
		return nil, models.ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, tx, pr.AuthorID)

	if err != nil {
		return nil, err
	}

	teamID := author.TeamName

	// Найти до двух рандомных членов команды
	reviewers, err := s.userRepo.GetRandomReviewers(ctx, tx, teamID, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	pr.Status = models.PROpen
	pr.Reviewers = reviewers
	pr.CreatedAt = time.Now().UTC()

	// Добавить, если такие есть в pull_request_reviewers и pull request
	if err := s.prRepo.Create(ctx, tx, pr); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return pr, nil
}
