package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandroag/desafio/app/gateway/db/postgres/migrations"
)

func New(addr string, minConn, maxConn int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing pgxpool config: %w", err)
	}

	// The defaults are located on top of pgxpool.pool.go
	config.MaxConns = maxConn
	config.MinConns = minConn

	pgxConn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("creating new pgxpool: %w", err)
	}

	if err = RunMigrations(addr, Migrations{
		Folder: ".",
		FS:     migrations.FS,
	}); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return pgxConn, nil
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}
