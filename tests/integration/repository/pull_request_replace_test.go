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

	"github.com/stretchr/testify/require"
)

func TestPullRequestRepository_ReplaceReviewer(t *testing.T) {
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

	// Создаём команду
	team, err := teamRepo.CreateTeam(ctx, db, "backend")
	require.NoError(t, err)

	// Создаём пользователей
	users := []models.User{
		{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{ID: "u3", Username: "Carl", TeamName: "backend", IsActive: true},
		{ID: "u4", Username: "Dima", TeamName: "backend", IsActive: true}, // кандидат на замену
	}

	teamModel := models.Team{
		ID:      team.ID,
		Name:    team.Name,
		Members: users,
	}

	_, err = userRepo.UpsertUsers(ctx, db, &teamModel)
	require.NoError(t, err)

	// Создаём PR
	pr := models.PullRequest{
		ID:        "pr-1",
		Name:      "Add Search",
		AuthorID:  "u1",
		Status:    "OPEN",
		Reviewers: []string{"u2", "u3"},
	}

	err = prRepo.Create(ctx, db, &pr)
	require.NoError(t, err)

	// Делаем замену ревьювера: u2 → u4
	err = prRepo.ReplaceReviewer(ctx, db, "pr-1", "u2", "u4")
	require.NoError(t, err)

	// Проверяем результат
	updated, err := prRepo.GetByID(ctx, db, "pr-1")
	require.NoError(t, err)

	require.ElementsMatch(t,
		[]string{"u4", "u3"},
		updated.Reviewers,
	)
}
