package pull_request_repository

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *prRepository) Create(ctx context.Context, q repositories.QueryExecer, pr *models.PullRequest) error {
	// Вставляем PR
	_, err := q.Exec(ctx,
		`INSERT INTO pull_requests (id, name, author_id, status, created_at)
         VALUES ($1, $2, $3, $4, $5)`,
		pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Добавляем reviewers
	batch := &pgx.Batch{}
	for _, uid := range pr.Reviewers {
		batch.Queue(`
            INSERT INTO pull_request_reviewers (pr_id, reviewer_id, created_at)
            VALUES ($1, $2, $3)
        `, pr.ID, uid, pr.CreatedAt)
	}

	br := q.SendBatch(ctx, batch)
	return br.Close()
}
