package prserv

import (
	"avito-reviewer/internal/models"
	"context"
	"math/rand"
)

//nolint:gocyclo
func (s *prService) ReassignReviewer(ctx context.Context, reasign *models.ReasignPR) (*models.PullRequest, string, error) {

	tx, err := s.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Проверить, что старый ревьювер существует
	hasUser, err := s.userRepo.UserExists(ctx, tx, reasign.OldReviewerID)
	if err != nil {
		return nil, "", err
	}

	if !hasUser {
		return nil, "", models.ErrNotFound
	}

	// Проверить существует ли pull_request
	hasPR, err := s.prRepo.Exists(ctx, tx, reasign.PRID)
	if err != nil {
		return nil, "", err
	}
	if !hasPR {
		return nil, "", models.ErrNotFound
	}

	// Загружаем PR полностью
	pr, err := s.prRepo.GetByID(ctx, tx, reasign.PRID)

	if err != nil {
		return nil, "", err
	}

	if pr.Status == models.PRMerged {
		return nil, "", models.ErrPRMerged
	}

	// Проверка, был ли старый ревьювер действительно назначен
	assigned := false
	for _, r := range pr.Reviewers {
		if r == reasign.OldReviewerID {
			assigned = true
			break
		}
	}

	if !assigned {
		return nil, "", models.ErrNotAssigned
	}

	//Получаем автора PR, чтобы определить команду
	author, err := s.userRepo.GetByID(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, "", models.ErrNotFound
	}

	// Найти кандидатов -> активные, из команды автора, не автор, не oldUser, не другие назначенные
	candidates, err := s.userRepo.GetReassignCandidates(
		ctx, tx, author.TeamName, pr.Reviewers, reasign.OldReviewerID, pr.AuthorID,
	)
	if err != nil {
		return nil, "", err
	}

	if len(candidates) == 0 {
		return nil, "", models.ErrNoCandidate
	}

	//nolint:gosec
	newReviewer := candidates[rand.Intn(len(candidates))] // Выбираем одного случайно

	for i := range pr.Reviewers {
		if pr.Reviewers[i] == reasign.OldReviewerID {
			pr.Reviewers[i] = newReviewer
			break
		}
	}

	if err := s.prRepo.ReplaceReviewer(ctx, tx, reasign.PRID, reasign.OldReviewerID, newReviewer); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return pr, newReviewer, nil
}
