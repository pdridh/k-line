package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/db"
	"github.com/pkg/errors"
)

type Store interface {
	CreateUser(ctx context.Context, name string, userType UserType, password string) (*User, error)
}

func NewPSQLStore(db *sqlx.DB) *sqlxStore {
	return &sqlxStore{
		db: db,
	}
}

type sqlxStore struct {
	db *sqlx.DB
}

func (s *sqlxStore) CreateUser(ctx context.Context, name string, userType UserType, password string) (*User, error) {
	q, a, err := db.PSQL.Insert("users").Columns("name", "type", "password").Values(name, userType, password).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var u User

	if err := s.db.QueryRowxContext(ctx, q, a...).StructScan(&u); err != nil {
		return nil, errors.Wrap(err, "failed to scan user into struct")
	}

	return &u, nil
}
