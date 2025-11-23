package prserv

import (
	"avito-reviewer/internal/models"
	"context"
	"time"
)

func (s *prService) MergePR(ctx context.Context, id string) (*models.PullRequest, error) {

	tx, err := s.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Загружаем PR со всеми ревьюерами
	pr, err := s.prRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	// Идемпотентность
	if pr.Status == models.PRMerged {
		_ = tx.Commit(ctx)
		return pr, nil
	}

	now := time.Now().UTC()
	pr.Status = models.PRMerged
	pr.MergedAt = &now

	// Меняем статус на мерж
	if err := s.prRepo.SetStatusOnMerged(ctx, tx, id, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return pr, nil
}
