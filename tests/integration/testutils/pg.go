package testutils

import (
	"avito-reviewer/internal/config"
	"avito-reviewer/pkg/postgres"
	"context"
	"time"

	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDB struct {
	Container *tcpg.PostgresContainer
	Cfg       *config.Config
}

//nolint:gosec
func StartPostgres(ctx context.Context) (*TestDB, error) {
	pg, err := tcpg.Run(ctx,
		"postgres:17",
		tcpg.WithDatabase("testdb"),
		tcpg.WithUsername("postgres"),
		tcpg.WithPassword("postgres"),
	)
	if err != nil {
		return nil, err
	}

	// даём контейнеру полностью подняться
	time.Sleep(2 * time.Second)

	host, err := pg.Host(ctx)
	if err != nil {
		return nil, err
	}

	mapped, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	cfg := &config.Config{
		Host:     host,
		Port:     uint16(mapped.Int()),
		Username: "postgres",
		Password: "postgres",
		Database: "testdb",
		MaxConns: 10,
		MinConns: 1,
	}

	// пытаемся подключиться + миграции
	for i := 0; i < 10; i++ {
		_, err = postgres.NewTest(ctx, cfg, "file://../../../db/migrations")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return &TestDB{Container: pg, Cfg: cfg}, err
}
