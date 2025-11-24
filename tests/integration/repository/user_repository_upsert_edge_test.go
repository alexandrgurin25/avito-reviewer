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

func TestUserRepository_UpsertUsers_EmptyMembers(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	team, err := teamRepo.CreateTeam(ctx, db, "t1")
	require.NoError(t, err)

	teamModel := &models.Team{
		ID:      team.ID,
		Name:    "t1",
		Members: []models.User{}, // пустой список
	}

	created, err := repo.UpsertUsers(ctx, db, teamModel)
	require.NoError(t, err)
	require.Len(t, created.Members, 0)
}

func TestUserRepository_UpsertUsers_UpdateExisting(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	team, err := teamRepo.CreateTeam(ctx, db, "t2")
	require.NoError(t, err)

	teamModel := &models.Team{
		ID:   team.ID,
		Name: team.Name,
		Members: []models.User{
			{ID: "u1", Username: "User1", IsActive: true},
		},
	}

	_, err = repo.UpsertUsers(ctx, db, teamModel)
	require.NoError(t, err)

	teamModel2 := &models.Team{
		ID:   team.ID,
		Name: team.Name,
		Members: []models.User{
			{ID: "u1", Username: "Updated", IsActive: false},
		},
	}

	updated, err := repo.UpsertUsers(ctx, db, teamModel2)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Members[0].Username)
	require.False(t, updated.Members[0].IsActive)
}

func TestUserRepository_UpsertUsers_TeamNotExists(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)

	teamModel := &models.Team{
		ID:   9999, // команды нет
		Name: "ghost",
		Members: []models.User{
			{ID: "u1", Username: "User", IsActive: true},
		},
	}

	_, err = repo.UpsertUsers(ctx, db, teamModel)
	require.Error(t, err)
	require.Contains(t, err.Error(), "foreign key")
}
