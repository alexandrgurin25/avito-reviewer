package postgres

import (
	"avito-reviewer/internal/config"
	"context"
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.MaxConns,
		cfg.MinConns,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {

		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	migrateConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	m, err := migrate.New(
		"file://./db/migrations",
		migrateConnString,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return pool, nil
}
