package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v4"
)

type Database struct {
	pool databaseInterface
}

type databaseInterface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

func New(dsn string) (*Database, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	db := &Database{
		pool: pool,
	}
	return db, nil
}

func (d *Database) Close() {
	d.pool.Close()
}

func NewMock() (*Database, pgxmock.PgxPoolIface, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		return nil, nil, err
	}
	db := &Database{
		pool: mock,
	}
	return db, mock, nil
}
