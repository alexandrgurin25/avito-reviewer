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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_UpsertAndGet(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	// Создаем команду
	team, err := teamRepo.CreateTeam(ctx, db, "user-test-team")
	require.NoError(t, err)

	// Тестовые пользователи
	testUsers := []models.User{
		{ID: "user-1", Username: "User One", IsActive: true},
		{ID: "user-2", Username: "User Two", IsActive: true},
		{ID: "user-3", Username: "User Three", IsActive: false},
	}

	teamModel := &models.Team{
		ID:      team.ID,
		Name:    team.Name,
		Members: testUsers,
	}

	// Тест создания пользователей
	createdTeam, err := repo.UpsertUsers(ctx, db, teamModel)
	require.NoError(t, err)
	assert.Len(t, createdTeam.Members, 3)

	// Тест получения пользователей команды
	users, err := repo.GetUsersByTeam(ctx, db, team.ID)
	require.NoError(t, err)
	assert.Len(t, users.Members, 3)

	// Тест изменения активности пользователя
	updatedUser, teamID, err := repo.SetActive(ctx, db, "user-3", true)
	require.NoError(t, err)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, team.ID, teamID)

	// Тест получения пользователя по ID
	user, err := repo.GetByID(ctx, db, "user-1")
	require.NoError(t, err)
	assert.Equal(t, "user-1", user.ID)
	assert.Equal(t, "User One", user.Username)
}
