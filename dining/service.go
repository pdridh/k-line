package dining

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/db/sqlc"
	"github.com/pkg/errors"
)

type service struct {
	Validate *validator.Validate
	store    db.Store
}

func NewService(v *validator.Validate, s db.Store) *service {
	return &service{
		Validate: v,
		store:    s,
	}
}

func (s *service) CreateOrder(ctx context.Context, tableID pgtype.Text, employeeID pgtype.UUID) (*pgtype.UUID, error) {

	t, err := s.store.GetTableByID(ctx, tableID.String)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, errors.Wrap(ErrUnknownTable, "store")
		}
		return nil, errors.Wrap(err, "store")
	}

	// Check if the table is available
	if t.Status != sqlc.TableStatusAvailable {
		return nil, errors.Wrap(ErrTableNotAvaliable, "store")
	}

	return s.store.CreateDiningOrderTx(ctx, tableID, employeeID)
}

func (s *service) CreateSession(ctx context.Context, tableID int) (*Session, error) {
	return nil, nil
}

func (s *service) IsTableAvailable(ctx context.Context, tableID int) (bool, error) {

	return false, nil
}

func (s *service) GetOngoingTable(ctx context.Context, tableID int) (*Session, error) {
	return nil, nil
}

func (s *service) AddItemsToSession(ctx context.Context, tableID int, items []SessionItem) ([]SessionItem, error) {
	return nil, nil
}

func (s *service) GetSessionItemsWithStatus(ctx context.Context, sessionID string, itemStatus ItemStatus) ([]SessionItem, error) {
	return nil, nil
}
