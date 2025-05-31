package menu

import (
	"strings"
	"time"
)

type MenuItem struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type MenuFilters struct {
	Search    string `json:"search"`
	Page      int    `json:"page"` // The request is sent as page but converted to offset for db
	Limit     int    `json:"limit"`
	OrderBy   string `json:"order_by"`
	SortOrder string `json:"sort_order"`
}

func (f *MenuFilters) Validate(allowedOrderBy map[string]bool, maxLimit int, defaultLimit int, defaultOrderBy string) {
	// Normalize pagination
	if f.Limit <= 0 || f.Limit > maxLimit {
		f.Limit = defaultLimit
	}

	if f.Page < 1 {
		f.Page = 1
	}

	// Normalize ordering
	if !allowedOrderBy[f.OrderBy] {
		f.OrderBy = defaultOrderBy
	}

	f.SortOrder = strings.ToUpper(f.SortOrder)
	if f.SortOrder != "ASC" && f.SortOrder != "DESC" {
		f.SortOrder = "ASC"
	}

	if f.Search != "" {
		f.Search = "%" + f.Search + "%"
	}
}
