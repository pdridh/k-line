package menu

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateItem(ctx context.Context, name string, description string, price float64) (*MenuItem, error)
	GetAllItems(ctx context.Context) ([]MenuItem, error)
	GetItemById(ctx context.Context, id int) (*MenuItem, error)
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
	const query = `INSERT INTO menu_items (name, description, price) VALUES(:name, :description, :price)	RETURNING *;`

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

func (s *sqlxStore) GetAllItems(ctx context.Context) ([]MenuItem, error) {
	const query = `SELECT * FROM menu_items;`

	// TODO add filtering and pagination probably
	var items []MenuItem

	if err := s.db.SelectContext(ctx, &items, query); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *sqlxStore) GetItemById(ctx context.Context, id int) (*MenuItem, error) {
	const query = `SELECT * FROM menu_items WHERE id = $1`

	var item MenuItem

	row := s.db.QueryRowxContext(ctx, query, id)
	if err := row.StructScan(&item); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}
