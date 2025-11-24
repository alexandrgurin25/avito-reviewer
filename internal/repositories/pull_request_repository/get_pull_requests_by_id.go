package prrepo

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"context"
)

func (*prRepository) GetPRsByID(ctx context.Context,
	db repositories.QueryExecer, ids []string) ([]models.PullRequest, error) {

	rows, err := db.Query(ctx,
		`SELECT id, name, author_id, status FROM pull_requests WHERE id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.PullRequest

	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		res = append(res, pr)
	}
	return res, nil
}
