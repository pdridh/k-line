package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pdridh/k-line/db/sqlc"
)

type Store interface {
	sqlc.Querier
}

type psqlStore struct {
	pool *pgxpool.Pool
	*sqlc.Queries
}

func NewPSQLStore(pool *pgxpool.Pool) Store {
	return &psqlStore{
		pool:    pool,
		Queries: sqlc.New(pool),
	}
}
