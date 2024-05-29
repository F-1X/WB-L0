package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB interface {
	Close()
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (Transaction, error)
}

type PostgresClientWrapper struct {
	pool *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, connString string) (PostgresDB, error) {
	pool, err := initPostgresClient(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &PostgresClientWrapper{pool: pool}, nil
}

func initPostgresClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

func (w *PostgresClientWrapper) Close() {
	w.pool.Close()
}

func (w *PostgresClientWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return w.pool.QueryRow(ctx, query, args...)
}

func (w *PostgresClientWrapper) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return w.pool.Query(ctx, query, args...)
}

func (w *PostgresClientWrapper) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return w.pool.Exec(ctx, query, args...)
}

func (w *PostgresClientWrapper) Begin(ctx context.Context) (Transaction, error) {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return nil, ErrBeginTransaction
	}
	return &PgTransaction{tx: tx}, nil
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)
	Exec(ctx context.Context, name string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type PgTransaction struct {
	tx pgx.Tx
}

func (t *PgTransaction) Commit(ctx context.Context) error {
	err := t.tx.Commit(ctx)
	if err != nil {
		return ErrCommitTransaction
	}
	return nil
}

func (t *PgTransaction) Rollback(ctx context.Context) error {
	err := t.tx.Rollback(ctx)
	if err != nil {
		return ErrRollbackTransaction
	}
	return nil
}

func (t *PgTransaction) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return t.tx.Prepare(ctx, name, sql)
}

func (t *PgTransaction) Exec(ctx context.Context, name string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, name, arguments...)
}
