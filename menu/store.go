package menu

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/api"
)

type Store interface {
	CreateItem(ctx context.Context, name string, description string, price float64) (*MenuItem, error)
	GetAllItems(ctx context.Context, filters MenuFilterArgs) ([]MenuItem, *api.PaginationMeta, error)
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

func (s *sqlxStore) GetAllItems(ctx context.Context, filters MenuFilterArgs) ([]MenuItem, *api.PaginationMeta, error) {
	countQuery := "SELECT COUNT(*) FROM menu_items"
	baseQuery := `SELECT * FROM menu_items`

	var total int
	meta := api.PaginationMeta{}

	if filters.Search != "" {
		baseQuery += ` WHERE name ILIKE :search`
		countQuery += ` WHERE name ILIKE :search`
		filters.Search = "%" + filters.Search + "%"
	}

	stmt, err := s.db.PrepareNamedContext(ctx, countQuery)
	if err != nil {
		return nil, nil, err
	}

	// Get total count
	if err := stmt.GetContext(ctx, &total, filters); err != nil {
		return nil, nil, err
	}

	meta.Total = total
	meta.TotalPages = (total + filters.Limit - 1) / filters.Limit
	if meta.Total == 0 {
		return nil, nil, nil
	}

	// Apply caps to limit
	if filters.Limit <= 0 || filters.Limit >= 50 {
		filters.Limit = 20
	}

	if filters.Page < 0 {
		filters.Page = 0
	}

	if filters.Page >= meta.TotalPages {
		filters.Page = meta.TotalPages
	}

	meta.Page = filters.Page
	meta.Limit = filters.Limit

	filters.Page = (filters.Page - 1) * filters.Limit

	var allowedOrderBy = map[string]bool{
		"created_at": true,
		"name":       true,
		"id":         true,
	}

	if !allowedOrderBy[filters.OrderBy] {
		filters.OrderBy = "name"
	}

	filters.SortOrder = strings.ToUpper(filters.SortOrder)
	if filters.SortOrder != "ASC" && filters.SortOrder != "DESC" {
		filters.SortOrder = "ASC"
	}

	baseQuery += fmt.Sprintf(` ORDER BY %s %s LIMIT :limit OFFSET :offset`, filters.OrderBy, filters.SortOrder)
	log.Println(baseQuery)

	stmt, err = s.db.PrepareNamedContext(ctx, baseQuery)
	if err != nil {
		return nil, nil, err
	}

	var items []MenuItem

	if err := stmt.SelectContext(ctx, &items, filters); err != nil {
		return nil, nil, err
	}

	return items, &meta, nil
}

func (s *sqlxStore) GetItemById(ctx context.Context, id int) (*MenuItem, error) {
	query := `SELECT * FROM menu_items WHERE id = ?`
	query = s.db.Rebind(query)

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
