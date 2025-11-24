package prrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
	"time"
)

func (*prRepository) SetStatusOnMerged(ctx context.Context, q repositories.QueryExecer, id string, mergedAt time.Time) error {
	tag, err := q.Exec(ctx,
		`UPDATE pull_requests
         SET status = 'MERGED', merged_at = $1
         WHERE id = $2`,
		mergedAt, id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
