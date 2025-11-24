package repository_test

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	"strconv"

	teamrepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"

	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetRandomReviewers(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)
	repo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)

	// Создаем команду с несколькими пользователями
	team, err := teamRepo.CreateTeam(ctx, db, "reviewers-team")
	require.NoError(t, err)

	users := []models.User{
		{ID: "author", Username: "Author", IsActive: true},
		{ID: "reviewer1", Username: "Reviewer 1", IsActive: true},
		{ID: "reviewer2", Username: "Reviewer 2", IsActive: true},
		{ID: "reviewer3", Username: "Reviewer 3", IsActive: false}, // неактивный
		{ID: "reviewer4", Username: "Reviewer 4", IsActive: true},
	}

	teamModel := &models.Team{
		ID:      team.ID,
		Name:    team.Name,
		Members: users,
	}

	_, err = repo.UpsertUsers(ctx, db, teamModel)
	require.NoError(t, err)

	teamID := strconv.Itoa(team.ID)

	// Тест получения случайных ревьюверов (исключая автора)
	reviewers, err := repo.GetRandomReviewers(ctx, db, teamID, "author")
	require.NoError(t, err)

	// Должны вернуться только активные пользователи, кроме автора
	assert.Len(t, reviewers, 2) // reviewer1, reviewer2, reviewer4 (3 активных минус автор = 2)

	// Проверяем, что автор не включен в ревьюверы
	for _, reviewer := range reviewers {
		assert.NotEqual(t, "author", reviewer)
	}
}
