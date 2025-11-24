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

func TestTeamAdd_AlreadyExists(t *testing.T) {
	ctx := context.Background()

	env, _ := testutils.StartPostgres(ctx)
	defer env.Container.Terminate(ctx)

	pool, _ := postgres.NewTest(ctx, env.Cfg, "file://../../../db/migrations")

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

	body := `{
		"team_name":"dup",
		"members":[{"user_id":"u1","username":"A"}]
	}`

	resp1, _ := http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(body))
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	resp2, _ := http.Post(ts.URL+"/team/add", "application/json", strings.NewReader(body))
	require.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	data, _ := io.ReadAll(resp2.Body)
	require.Contains(t, string(data), "TEAM_EXISTS")
}
