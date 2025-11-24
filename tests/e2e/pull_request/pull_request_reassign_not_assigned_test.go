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

func TestReassignNotAssigned(t *testing.T) {
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

	userSrv := userserv.NewService(userRepo, teamRepo, prRepo, db)
	teamSrv := teamserv.NewService(userRepo, teamRepo, db)
	prSrv := prserv.NewService(userRepo, prRepo, db)

	router := build.Router(
		teamhand.NewTeamHandler(teamSrv),
		prhand.NewPRHandler(prSrv),
		userhand.NewTeamHandler(userSrv),
	)

	ts := httptest.NewServer(router)
	defer ts.Close()

	team1 := `{
		"team_name": "t1",
		"members":[
			{"user_id":"u1","username":"A","is_active":true},
			{"user_id":"u2","username":"B","is_active":true},
			{"user_id":"u3","username":"C","is_active":true}
		]
	}`
	http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(team1))

	team2 := `{
		"team_name": "t2",
		"members":[
			{"user_id":"u4","username":"A","is_active":true},
			{"user_id":"u5","username":"B","is_active":true},
			{"user_id":"u6","username":"C","is_active":true}
		]
	}`

	http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(team2))

	pr := `{
		"pull_request_id":"pr-x",
		"pull_request_name":"Test",
		"author_id":"u1"
	}`
	http.Post(ts.URL+"/pullRequest/create", "application/json", strings.NewReader(pr))

	body := `{"pull_request_id":"pr-x","old_user_id":"u6"}`

	resp, _ := http.Post(ts.URL+"/pullRequest/reassign", "application/json", strings.NewReader(body))
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(data), `"code":"NOT_ASSIGNED"`)
}
