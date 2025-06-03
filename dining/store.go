package dining

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/db"
	"github.com/pkg/errors"
)

type Store interface {
	CreateSession(ctx context.Context, tableID int) (*Session, error)
	GetOngoingSessionByTable(ctx context.Context, tableID int) (*Session, error)
	CreateSessionItems(ctx context.Context, sessionID uuid.UUID, items []SessionItem) ([]SessionItem, error)
	GetSessionItemsWithStatus(ctx context.Context, sessionID uuid.UUID, itemStatus ItemStatus) ([]SessionItem, error)
}

func NewPSQLStore(db *sqlx.DB) *sqlxStore {
	return &sqlxStore{
		db: db,
	}
}

type sqlxStore struct {
	db *sqlx.DB
}

func (s *sqlxStore) GetOngoingSessionByTable(ctx context.Context, tableID int) (*Session, error) {
	q, a, err := db.PSQL.Select("*").From("dining_sessions").Where("table_id = ? AND status = 'ongoing'", tableID).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var session Session

	if err := s.db.QueryRowxContext(ctx, q, a...).StructScan(&session); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "scan")
	}

	return &session, nil
}

func (s *sqlxStore) CreateSession(ctx context.Context, tableID int) (*Session, error) {
	q, a, err := db.PSQL.Insert("dining_sessions").Columns("table_id").Values(tableID).Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var session Session

	if err := s.db.QueryRowxContext(ctx, q, a...).StructScan(&session); err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return &session, nil
}

func (s *sqlxStore) CreateSessionItems(ctx context.Context, sessionID uuid.UUID, items []SessionItem) ([]SessionItem, error) {

	baseQuery := db.PSQL.Insert("dining_items").Columns("session_id", "item_id", "quantity")
	for _, item := range items {
		baseQuery = baseQuery.Values(sessionID, item.ItemID, item.Quantity)
	}
	q, a, err := baseQuery.Suffix("RETURNING *").ToSql()
	if err != nil {
		return []SessionItem{}, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Queryx(q, a...)
	if err != nil {
		return []SessionItem{}, err
	}
	defer rows.Close()

	var inserted []SessionItem
	for rows.Next() {
		var i SessionItem
		if err := rows.StructScan(&i); err != nil {
			return []SessionItem{}, err
		}
		inserted = append(inserted, i)
	}

	return inserted, nil
}

func (s *sqlxStore) GetSessionItemsWithStatus(ctx context.Context, sessionID uuid.UUID, itemStatus ItemStatus) ([]SessionItem, error) {
	q, a, err := db.PSQL.Select("*").From("dining_items").Where("session_id = ? AND status = ?", sessionID, itemStatus).ToSql()
	if err != nil {
		return []SessionItem{}, errors.Wrap(err, "failed to build query")
	}

	items := []SessionItem{}
	if err := s.db.SelectContext(ctx, &items, q, a...); err != nil {
		return []SessionItem{}, errors.Wrap(err, "failed to get items")
	}

	return items, nil
}
