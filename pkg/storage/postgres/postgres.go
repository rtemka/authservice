package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres handles CRUD operation over the
// postgresql database
type Postgres struct {
	db *pgxpool.Pool
}

// New returns new *Postgres object
func New(connstr string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	pgres := &Postgres{
		db: pool,
	}

	return pgres, pgres.db.Ping(context.Background())
}
