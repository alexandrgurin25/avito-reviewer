package repository_test

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	prrepo "avito-reviewer/internal/repositories/pull_request_repository"
	teamrepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"
	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPullRequestRepository_Merge(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	prRepo := prrepo.NewPRRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	// team + users
	team, err := teamRepo.CreateTeam(ctx, db, "backend")
	require.NoError(t, err)

	users := []models.User{
		{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
	}
	teamModel := models.Team{
		ID:      team.ID,
		Name:    team.Name,
		Members: users,
	}
	_, err = userRepo.UpsertUsers(ctx, db, &teamModel)
	require.NoError(t, err)

	// PR
	pr := models.PullRequest{
		ID:        "pr-m1",
		Name:      "Feature",
		AuthorID:  "u1",
		Status:    "OPEN",
		Reviewers: []string{"u2"},
	}
	err = prRepo.Create(ctx, db, &pr)
	require.NoError(t, err)

	// merge
	err = prRepo.SetStatusOnMerged(ctx, db, "pr-m1", time.Now().UTC())
	require.NoError(t, err)

	// check
	merged, err := prRepo.GetByID(ctx, db, "pr-m1")
	require.NoError(t, err)

	require.Equal(t, "MERGED", string(merged.Status))
	require.NotNil(t, merged.MergedAt)
}
