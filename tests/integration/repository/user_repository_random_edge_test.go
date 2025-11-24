package repository_test

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	teamrepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"
	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetRandomReviewers_NoActiveUsers(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	team, _ := teamRepo.CreateTeam(ctx, db, "test")

	users := []models.User{
		{ID: "u1", Username: "A", IsActive: false},
		{ID: "u2", Username: "B", IsActive: false},
	}

	teamModel := &models.Team{ID: team.ID, Name: team.Name, Members: users}
	repo.UpsertUsers(ctx, db, teamModel)

	revs, err := repo.GetRandomReviewers(ctx, db, "author", "u1")
	require.NoError(t, err)
	require.Len(t, revs, 0)
}

func TestUserRepository_GetRandomReviewers_AuthorOnlyActive(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	team, _ := teamRepo.CreateTeam(ctx, db, "test")

	users := []models.User{
		{ID: "u1", Username: "Author", IsActive: true},
		{ID: "u2", Username: "Inactive", IsActive: false},
	}

	teamModel := &models.Team{ID: team.ID, Name: team.Name, Members: users}
	repo.UpsertUsers(ctx, db, teamModel)

	revs, err := repo.GetRandomReviewers(ctx, db, "test", "u1")
	require.NoError(t, err)
	require.Len(t, revs, 0)
}

func TestUserRepository_GetRandomReviewers_TeamNotFound(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)

	revs, err := repo.GetRandomReviewers(ctx, db, "unknown-team", "u1")
	require.NoError(t, err)
	require.Len(t, revs, 0)
}
