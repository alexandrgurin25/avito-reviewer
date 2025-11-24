package e2e_test

import (
	build "avito-reviewer/internal/handlers/build_router"
	prhand "avito-reviewer/internal/handlers/pullRequest"
	teamhand "avito-reviewer/internal/handlers/team_handler"
	userhand "avito-reviewer/internal/handlers/user_handler"
	"avito-reviewer/internal/repositories"
	prrepo "avito-reviewer/internal/repositories/pull_request_repository"
	teamrepo "avito-reviewer/internal/repositories/team_repository"
	userrepo "avito-reviewer/internal/repositories/user_repository"
	prserv "avito-reviewer/internal/services/pull_request_services"
	teamserv "avito-reviewer/internal/services/team_services"
	userserv "avito-reviewer/internal/services/user_services"
	"avito-reviewer/pkg/postgres"
	"avito-reviewer/tests/integration/testutils"

	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePullRequest(t *testing.T) {
	ctx := context.Background()

	// запуск тестовой базы
	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	repository := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(repository)
	teamRepo := teamrepo.NewTeamRepository(repository)
	prRepo := prrepo.NewPRRepository(repository)

	userSrv := userserv.NewService(userRepo, teamRepo, prRepo, repository)
	teamSrv := teamserv.NewService(userRepo, teamRepo, repository)
	prSrv := prserv.NewService(userRepo, prRepo, repository)

	teamHandler := teamhand.NewTeamHandler(teamSrv)
	prHandler := prhand.NewPRHandler(prSrv)
	userHandler := userhand.NewTeamHandler(userSrv)

	router := build.Router(teamHandler, prHandler, userHandler)

	ts := httptest.NewServer(router)
	defer ts.Close()

	// создаём команду
	team := `{
		"team_name": "backend",
		"members":[
			{"user_id":"u1","username":"A","is_active":true},
			{"user_id":"u2","username":"B","is_active":true}
		]
	}`
	_, err = http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(team))
	require.NoError(t, err)

	// создаём PR
	pr := `{
		"pull_request_id":"pr-1",
		"pull_request_name":"Add search",
		"author_id":"u1"
	}`
	resp, err := http.Post(ts.URL+"/pullRequest/create", "application/json", strings.NewReader(pr))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}
