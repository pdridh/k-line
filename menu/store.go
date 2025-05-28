package menu

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateItem(ctx context.Context, name string, description string, price float64) (*MenuItem, error)
}

func NewPSQLStore(db *sqlx.DB) *sqlxStore {
	return &sqlxStore{
		db: db,
	}
}

type sqlxStore struct {
	db *sqlx.DB
}

// Insert a menu item using the given name, description and price and returns the MenuItem filled with all the fields
// in the table.
func (s *sqlxStore) CreateItem(ctx context.Context, name string, description string, price float64) (*MenuItem, error) {
	const query = `INSERT INTO menu_items (name, description, price) VALUES(:name, :description, :price)	RETURNING *`

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var i MenuItem

	namedArgs := map[string]any{
		"name":        name,
		"description": description,
		"price":       price,
	}

	if err := stmt.GetContext(ctx, &i, namedArgs); err != nil {
		return nil, err
	}

	return &i, nil
}
