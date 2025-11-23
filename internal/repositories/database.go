package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mocks/row_mock.go -package=mocks github.com/jackc/pgx/v5 Row
type QueryExecer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type DB interface {
	QueryExecer
	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	QueryExecer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type PgxPoolAdapter struct {
	pool *pgxpool.Pool
}

func NewPgxPoolAdapter(pool *pgxpool.Pool) *PgxPoolAdapter {
	return &PgxPoolAdapter{pool: pool}
}

func (a *PgxPoolAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return a.pool.QueryRow(ctx, sql, args...)
}

func (a *PgxPoolAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return a.pool.Exec(ctx, sql, args...)
}

func (a *PgxPoolAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return a.pool.Query(ctx, sql, args...)
}

func (a *PgxPoolAdapter) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return a.pool.SendBatch(ctx, b)
}

func (a *PgxPoolAdapter) BeginTx(ctx context.Context) (Tx, error) {
	pgxTx, err := a.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &PgxTxAdapter{tx: pgxTx}, nil
}

type PgxTxAdapter struct {
	tx pgx.Tx
}

func (t *PgxTxAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *PgxTxAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *PgxTxAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, sql, args...)
}

func (t *PgxTxAdapter) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *PgxTxAdapter) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *PgxTxAdapter) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return t.tx.SendBatch(ctx, b)
}
