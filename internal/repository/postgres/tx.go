package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX matches sqlc's expected interface for pgx.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TxManager provides transaction boundary for usecases.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx DBTX) error) error
}

type TxManagerPGX struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManagerPGX {
	return &TxManagerPGX{pool: pool}
}

func (m *TxManagerPGX) WithinTx(ctx context.Context, fn func(ctx context.Context, tx DBTX) error) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
