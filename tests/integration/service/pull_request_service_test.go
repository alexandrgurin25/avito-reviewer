package service_test

import (
	"avito-reviewer/internal/models"
	"avito-reviewer/internal/repositories"
	prrepo "avito-reviewer/internal/repositories/pull_request_repository"
	teamrepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"

	prserv "avito-reviewer/internal/services/pull_request_services"
	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"

	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPRService_Create_OK(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)
	prRepo := prrepo.NewPRRepository(db)
	service := prserv.NewService(userRepo, prRepo, db)

	team, _ := teamRepo.CreateTeam(ctx, db, "backend")

	users := []models.User{
		{ID: "u1", Username: "Author", IsActive: true, TeamName: "backend"},
		{ID: "u2", Username: "Rev1", IsActive: true, TeamName: "backend"},
		{ID: "u3", Username: "Rev2", IsActive: true, TeamName: "backend"},
	}

	userRepo.UpsertUsers(ctx, db, &models.Team{
		ID: team.ID, Name: team.Name, Members: users,
	})

	pr, err := service.CreatePullRequest(ctx, &models.PullRequest{
		ID: "pr-x", Name: "Add Feature", AuthorID: "u1",
	})
	require.NoError(t, err)
	require.Len(t, pr.Reviewers, 2)
}

func TestPRService_Create_UserNotFound(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	prRepo := prrepo.NewPRRepository(db)

	service := prserv.NewService(userRepo, prRepo, db)

	_, err = service.CreatePullRequest(ctx, &models.PullRequest{
		ID: "pr-x", Name: "Test", AuthorID: "ghost",
	})
	require.Error(t, err)
}

func TestPRService_Merge_OK(t *testing.T) {
	ctx := context.Background()

	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)
	prRepo := prrepo.NewPRRepository(db)
	service := prserv.NewService(userRepo, prRepo, db)

	team, _ := teamRepo.CreateTeam(ctx, db, "backend")

	userRepo.UpsertUsers(ctx, db, &models.Team{
		ID: team.ID, Name: "backend",
		Members: []models.User{
			{ID: "u1", Username: "A", IsActive: true, TeamName: "backend"},
		},
	})

	pr := &models.PullRequest{
		ID: "pr-m", Name: "Test", AuthorID: "u1",
		Status: models.PROpen, Reviewers: []string{},
	}
	prRepo.Create(ctx, db, pr)

	out, err := service.MergePR(ctx, "pr-m")
	require.NoError(t, err)
	require.Equal(t, models.PRMerged, out.Status)
}

func TestPRService_Merge_NotFound(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	prRepo := prrepo.NewPRRepository(db)

	service := prserv.NewService(userRepo, prRepo, db)

	_, err := service.MergePR(ctx, "nope")
	require.Error(t, err)
}

func TestPRService_Reassign_OK(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)
	prRepo := prrepo.NewPRRepository(db)
	service := prserv.NewService(userRepo, prRepo, db)

	team, _ := teamRepo.CreateTeam(ctx, db, "backend")

	users := []models.User{
		{ID: "u1", Username: "Author", IsActive: true, TeamName: "backend"},
		{ID: "u2", Username: "Old", IsActive: true, TeamName: "backend"},
		{ID: "u3", Username: "New", IsActive: true, TeamName: "backend"},
	}
	userRepo.UpsertUsers(ctx, db, &models.Team{
		ID: team.ID, Name: "backend", Members: users,
	})

	prRepo.Create(ctx, db, &models.PullRequest{
		ID: "pr-r", Name: "Test", AuthorID: "u1",
		Status: models.PROpen, Reviewers: []string{"u2"},
	})

	out, newR, err := service.ReassignReviewer(ctx, &models.ReasignPR{
		PRID: "pr-r", OldReviewerID: "u2",
	})
	require.NoError(t, err)
	require.Equal(t, "u3", newR)
	require.Contains(t, out.Reviewers, "u3")
}

func TestPRService_Reassign_NoCandidate(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)
	prRepo := prrepo.NewPRRepository(db)
	service := prserv.NewService(userRepo, prRepo, db)

	team, _ := teamRepo.CreateTeam(ctx, db, "backend")

	users := []models.User{
		{ID: "u1", Username: "Author", IsActive: true},
		{ID: "u2", Username: "Old", IsActive: true},
		{ID: "u3", Username: "New", IsActive: false},
	}
	for i := range users {
		users[i].TeamName = "backend"
	}

	userRepo.UpsertUsers(ctx, db, &models.Team{
		ID: team.ID, Name: team.Name, Members: users,
	})

	prRepo.Create(ctx, db, &models.PullRequest{
		ID: "pr-no", Name: "X", AuthorID: "u1", Status: models.PROpen,
		Reviewers: []string{"u2"},
	})

	_, _, err := service.ReassignReviewer(ctx, &models.ReasignPR{
		PRID: "pr-no", OldReviewerID: "u2",
	})
	require.Error(t, err)
}
