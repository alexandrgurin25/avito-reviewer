package pull_request_repository

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (r *prRepository) GetIDByReviewerId(ctx context.Context, q repositories.QueryExecer, userID string) ([]string, error) {

	// Подтянем id ревью
	rows, err := q.Query(ctx,
		`SELECT pr_id 
         FROM pull_request_reviewers
         WHERE reviewer_id = $1`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prIDs []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		prIDs = append(prIDs, uid)
	}

	return prIDs, nil
}
