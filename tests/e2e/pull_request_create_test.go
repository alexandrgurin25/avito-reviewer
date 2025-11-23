package e2e_test

import (
	"avito-reviewer/internal/config"
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
	"time"

	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	p "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestCreatePullRequest(t *testing.T) {
	ctx := context.Background()

	// 1) Поднимаем PostgreSQL
	pg, err := p.Run(ctx,
		"postgres:17",
		p.WithDatabase("testdb"),
		p.WithUsername("postgres"),
		p.WithPassword("postgres"),
	)
	require.NoError(t, err)
	defer pg.Terminate(ctx)

	// Даем время PostgreSQL полностью запуститься
	time.Sleep(5 * time.Second)

	host, err := pg.Host(ctx)
	require.NoError(t, err)

	port, err := pg.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Преобразуем порт в uint16
	portUint16 := uint16(port.Int())

	cfg := &config.Config{
		Host:     host,
		Port:     portUint16,
		Username: "postgres",
		Password: "postgres",
		Database: "testdb",
		MaxConns: 10,
		MinConns: 1,
	}

	// Добавляем retry для подключения к БД
	var db *pgxpool.Pool
	for i := 0; i < 10; i++ {
		db, err = postgres.NewTest(ctx, cfg, "file://../../db/migrations")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	require.NoError(t, err)
	// Остальной код без изменений...
	repository := repositories.NewPgxPoolAdapter(db)

	userRepo := userrepo.NewUserRepository(repository)
	teamRepo := teamrepo.NewTeamRepository(repository)
	prRepo := prrepo.NewPRRepository(repository)

	userSrv := userserv.NewService(userRepo, teamRepo, prRepo, repository)
	teamSrv := teamserv.NewService(userRepo, teamRepo, repository)
	prSrv := prserv.NewService(userRepo, prRepo)

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
