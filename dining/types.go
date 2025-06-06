package dining

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SessionStatus string

const (
	SessionOngoing   SessionStatus = "ongoing"
	SessionCompleted SessionStatus = "completed"
	SessionCancelled SessionStatus = "cancelled"
)

type Session struct {
	ID          pgtype.UUID   `db:"id"`
	Status      SessionStatus `db:"status"`
	TableID     int           `db:"table_id"`
	StartedAt   time.Time     `db:"started_at"`
	CompletedAt *time.Time    `db:"completed_at"`
}

type ItemStatus string

const (
	ItemPending   = "pending"
	ItemPreparing = "preparing"
	ItemReady     = "ready"
	ItemCompleted = "completed"
	ItemCancelled = "cancelled"
)

type SessionItem struct {
	ID        int         `db:"id"`
	SessionID pgtype.UUID `db:"session_id"`
	ItemID    int         `db:"item_id" json:"item_id"`
	Quantity  int         `db:"quantity" json:"quantity"`
	Status    ItemStatus  `db:"status"`
}
