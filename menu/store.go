package menu

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db"
	"github.com/pkg/errors"
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
	query, args, err := db.PSQL.Insert("menu_items").Columns("name", "description", "price").Values(name, description, price).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var i MenuItem

	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&i); err != nil {
		return nil, errors.Wrap(err, "failed to scan item into MenuItem")
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
		return nil, nil, errors.Wrap(err, "failed to get count")
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

	queryString, args, err := finalQuery.ToSql()
	if err != nil {
		return []MenuItem{}, nil, errors.Wrap(err, "failed to build query")
	}

	items := []MenuItem{}
	if err := s.db.SelectContext(ctx, &items, queryString, args...); err != nil {
		return []MenuItem{}, nil, errors.Wrap(err, "failed to get items")
	}

	return items, meta, nil
}

func (s *sqlxStore) GetItemById(ctx context.Context, id int) (*MenuItem, error) {
	baseQuery := db.PSQL.Select("*").From("menu_items").Where("id = ?", id)

	queryString, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var item MenuItem

	row := s.db.QueryRowxContext(ctx, queryString, args...)
	if err := row.StructScan(&item); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetItemById failed to load row into struct")
	}

	return &item, nil
}
