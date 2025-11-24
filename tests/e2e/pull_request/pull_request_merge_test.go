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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergePullRequest(t *testing.T) {
	ctx := context.Background()

	// запуск тестовой базы
	env, err := testutils.StartPostgres(ctx)
	require.NoError(t, err)
	defer env.Container.Terminate(ctx)

	pool, err := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")
	require.NoError(t, err)

	// поднимаем слой приложения 
	db := repositories.NewPgxPoolAdapter(pool)

	userRepo := userrepo.NewUserRepository(db)
	teamRepo := teamrepo.NewTeamRepository(db)
	prRepo := prrepo.NewPRRepository(db)

	userSrv := userserv.NewService(userRepo, teamRepo, prRepo, db)
	teamSrv := teamserv.NewService(userRepo, teamRepo, db)
	prSrv := prserv.NewService(userRepo, prRepo, db)

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
			{"user_id":"u1","username":"Alice","is_active":true},
			{"user_id":"u2","username":"Bob","is_active":true}
		]
	}`
	_, err = http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(team))
	require.NoError(t, err)

	// создаём PR
	pr := `{
		"pull_request_id":"pr-100",
		"pull_request_name":"Add Search",
		"author_id":"u1"
	}`
	resp, err := http.Post(ts.URL+"/pullRequest/create", "application/json", strings.NewReader(pr))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// merge PR
	mergeBody := `{
		"pull_request_id":"pr-100"
	}`
	resp, err = http.Post(ts.URL+"/pullRequest/merge", "application/json", strings.NewReader(mergeBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// проверяем, что PR стал merged
	getResp, err := http.Post(ts.URL+"/pullRequest/merge", "application/json", strings.NewReader(mergeBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getResp.StatusCode)

	data, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)

	body := string(data)
	require.Contains(t, body, `"status":"MERGED"`)
	require.Contains(t, body, `"mergedAt"`)

}
