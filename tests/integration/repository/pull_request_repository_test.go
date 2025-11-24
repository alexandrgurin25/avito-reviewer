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

func TestPullRequestRepository_CRUD(t *testing.T) {
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

	// Создаём команду, чтобы у пользователя была ссылка
	createdTeam, err := teamRepo.CreateTeam(ctx, db, "backend")
	require.NoError(t, err)

	user1 := models.User{
		ID:       "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	user2 := models.User{
		ID:       "u2",
		Username: "Alex",
		TeamName: "backend",
		IsActive: true,
	}

	user3 := models.User{
		ID:       "u3",
		Username: "Ali",
		TeamName: "backend",
		IsActive: true,
	}

	team := models.Team{
		ID:      createdTeam.ID,
		Name:    createdTeam.Name,
		Members: []models.User{1: user1, 2: user2, 3: user3},
	}

	_, err = userRepo.UpsertUsers(ctx, db, &team)
	require.NoError(t, err)

	// PR ещё нет
	exists, err := prRepo.Exists(ctx, db, "pr-1")
	require.NoError(t, err)
	require.False(t, exists)

	// создаём PR
	createdPR := models.PullRequest{
		ID:        "pr-1",
		Name:      "Add Search",
		AuthorID:  "u1",
		Status:    "OPEN",
		Reviewers: []string{"u2", "u3"},
	}
	err = prRepo.Create(ctx, db, &createdPR)
	require.NoError(t, err)

	// читаем PR
	readPR, err := prRepo.GetByID(ctx, db, "pr-1")
	require.NoError(t, err)
	require.Equal(t, createdPR.AuthorID, readPR.AuthorID)
	require.Equal(t, createdPR.Reviewers, readPR.Reviewers)
}
