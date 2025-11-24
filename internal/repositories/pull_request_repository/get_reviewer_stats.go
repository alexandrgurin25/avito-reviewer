package prrepo

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (*prRepository) GetReviewerStats(ctx context.Context, db repositories.QueryExecer) (map[string]int, error) {
	rows, err := db.Query(ctx,
		`SELECT reviewer_id, COUNT(*) 
		 FROM pull_request_reviewers
		 GROUP BY reviewer_id`)
		 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]int{}

	for rows.Next() {
		var id string
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, err
		}
		stats[id] = count
	}

	return stats, nil
}
