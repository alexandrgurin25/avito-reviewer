package pull_request_repository

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *prRepository) ReplaceReviewer(ctx context.Context, q repositories.QueryExecer, prID, oldID, newID string) error {
	tag, err := q.Exec(ctx,
		`UPDATE pull_request_reviewers
         SET reviewer_id = $1
         WHERE pr_id = $2 AND reviewer_id = $3`,
		newID, prID, oldID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotAssigned
	}
	return nil
}
