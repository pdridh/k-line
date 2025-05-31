package menu

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db"
)

type Store interface {
	CreateItem(ctx context.Context, name string, description string, price float64) (*MenuItem, error)
	GetAllItems(ctx context.Context, filters *MenuFilters) ([]MenuItem, *api.PaginationMeta, error)
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

func (s *sqlxStore) GetAllItems(ctx context.Context, filters *MenuFilters) ([]MenuItem, *api.PaginationMeta, error) {
	baseQuery := db.PSQL.Select("*").From("menu_items")
	countQuery := db.PSQL.Select("COUNT(*)").From("menu_items")

	var total int

	// Validate filters
	allowedOrderBy := map[string]bool{
		"name":       true,
		"created_at": true,
		"id":         true,
	}
	filters.Validate(allowedOrderBy, 50, 20, "name")

	// Apply filters
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		baseQuery = baseQuery.Where(squirrel.ILike{"name": searchPattern})
		countQuery = countQuery.Where(squirrel.ILike{"name": searchPattern})
	}

	total, err := db.GetCount(ctx, s.db, countQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("store: failed to get count: %w", err)
	}

	if total == 0 {
		return []MenuItem{}, nil, nil
	}

	meta := api.CalculatePaginationMeta(total, filters.Page, filters.Limit)

	offset := (filters.Page - 1) * filters.Limit
	finalQuery := baseQuery.
		OrderBy(filters.OrderBy + " " + filters.SortOrder).
		Limit(uint64(filters.Limit)).
		Offset(uint64(offset))

	sql, args, err := finalQuery.ToSql()
	if err != nil {
		return []MenuItem{}, nil, fmt.Errorf("store: failed to build query: %w", err)
	}

	items := []MenuItem{}
	if err := s.db.SelectContext(ctx, &items, sql, args...); err != nil {
		return []MenuItem{}, nil, fmt.Errorf("store: failed to get items: %w", err)
	}

	return items, meta, nil
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
