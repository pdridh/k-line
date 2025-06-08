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

func (s *service) AddItemsToOrder(ctx context.Context, orderID pgtype.UUID, items []RequestItem) error {
	// Check if the orderID is open and shit
	o, err := s.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errors.Wrap(ErrUnknownOrder, "store")
		}
		return errors.Wrap(err, "store")
	}

	if o.Status != sqlc.OrderStatusOngoing {
		return errors.Wrap(ErrOrderNotOngoing, "store")
	}

	var arg sqlc.AddOrderItemsBulkParams
	arg.OrderID = o.ID

	for _, i := range items {
		arg.ItemIds = append(arg.ItemIds, int32(i.ItemID))
		arg.Quantity = append(arg.Quantity, int32(i.Quantity))
		arg.Notes = append(arg.Notes, i.Note)
	}

	if err := s.store.AddOrderItemsBulk(ctx, arg); err != nil {
		return err
	}

	return nil
}
