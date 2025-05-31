package user

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pdridh/k-line/db"
	"github.com/pkg/errors"
)

type Store interface {
	CreateUser(ctx context.Context, email string, name string, userType UserType, password string) (*User, error)
}

func NewPSQLStore(db *sqlx.DB) *sqlxStore {
	return &sqlxStore{
		db: db,
	}
}

type sqlxStore struct {
	db *sqlx.DB
}

func (s *sqlxStore) CreateUser(ctx context.Context, email string, name string, userType UserType, password string) (*User, error) {
	// TODO turn it into hashed passwords
	q, a, err := db.PSQL.Insert("users").Columns("email", "name", "type", "password").Values(email, name, userType, password).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var u User

	if err := s.db.QueryRowxContext(ctx, q, a...).StructScan(&u); err != nil {

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" && strings.Contains(pqErr.Constraint, "email") {
			return nil, errors.Wrap(ErrDuplicateEmail, "failed to create user")
		}

		return nil, errors.Wrap(err, "failed to scan user into struct")
	}

	return &u, nil
}
