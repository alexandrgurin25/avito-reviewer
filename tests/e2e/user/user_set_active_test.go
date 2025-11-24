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

func TestUserSetActive(t *testing.T) {
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

	// создаем команду
	team := `{
		"team_name": "t-active",
		"members":[
			{"user_id":"u1","username":"A","is_active":false}
		]
	}`
	http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(team))

	// активируем пользователя
	body := `{"user_id":"u1","is_active":true}`
	resp, err := http.Post(ts.URL+"/users/setIsActive", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(data), `"is_active":true`)
}
