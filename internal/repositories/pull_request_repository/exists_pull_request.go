package prrepo

import (
	"avito-reviewer/internal/repositories"
	"context"
)

func (*prRepository) Exists(ctx context.Context, q repositories.QueryExecer, id string) (bool, error) {
	var exists bool
	err := q.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id=$1)`,
		id,
	).Scan(&exists)

	return exists, err
}
