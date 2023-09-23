package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

type Tx interface {
	DB
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func New(ctx context.Context, url string) (*Postgres, error) {
	var pg Postgres
	var err error

	pg.Pool, err = pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	return &pg, nil
}

func (p *Postgres) Begin(ctx context.Context) (Tx, error) {
	return p.Pool.Begin(ctx)
}
