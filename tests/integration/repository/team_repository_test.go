package repository_test

import (
	"avito-reviewer/internal/repositories"
	teamrepo "avito-reviewer/internal/repositories/team_repository"
	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTeamRepository_CreateAndExists(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := teamrepo.NewTeamRepository(db)

	// команда ещё не существует
	exists, err := repo.TeamExists(ctx, db, "backend")
	require.NoError(t, err)
	require.False(t, exists)

	// создаём команду
	_, err = repo.CreateTeam(ctx, db, "backend")
	require.NoError(t, err)

	// теперь должна существовать
	exists, err = repo.TeamExists(ctx, db, "backend")
	require.NoError(t, err)
	require.True(t, exists)
}
