package dining

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/db/sqlc"
)

type RequestItem struct {
	ItemID   int    `json:"item_id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"notes"`
}

type Table struct {
	ID       string           `json:"id"`
	Capacity int16            `json:"capacity"`
	Status   sqlc.TableStatus `json:"status"`
	Notes    pgtype.Text      `json:"notes"`
}

type Order struct {
	ID          pgtype.UUID      `json:"id"`
	Type        sqlc.OrderType   `json:"type"`
	EmployeeID  pgtype.UUID      `json:"employee_id"`
	Status      sqlc.OrderStatus `json:"status"`
	TableID     pgtype.Text      `json:"table_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	CompletedAt pgtype.Timestamp `json:"completed_at"`
}
