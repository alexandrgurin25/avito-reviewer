package prrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"time"
)

func (*prRepository) GetByID(ctx context.Context, q repositories.QueryExecer, id string) (*models.PullRequest, error) {
	var pr models.PullRequest
	var mergedAt *time.Time

	err := q.QueryRow(ctx,
		`SELECT id, name, author_id, status, created_at, merged_at
         FROM pull_requests
         WHERE id = $1`,
		id,
	).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)

	if err != nil {
		return nil, err
	}

	pr.MergedAt = mergedAt

	// Подтянем ревьюеров
	rows, err := q.Query(ctx,
		`SELECT reviewer_id 
         FROM pull_request_reviewers
         WHERE pr_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, uid)
	}

	pr.Reviewers = reviewers
	return &pr, nil
}
