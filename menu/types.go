package menu

import "time"

type MenuItem struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type MenuFilterArgs struct {
	Search    string `db:"search" json:"search"`
	Page      int    `db:"offset" json:"page"` // The request is sent as page but converted to offset for db
	Limit     int    `db:"limit" json:"limit"`
	OrderBy   string `db:"order_by" json:"order_by"`
	SortOrder string `db:"sort_order" json:"sort_order"`
}
