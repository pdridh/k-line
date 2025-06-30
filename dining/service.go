package dining

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
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
			return nil, errors.Wrap(api.ErrUnknownTable.Error, "store")
		}
		return nil, errors.Wrap(err, "store")
	}

	// Check if the table is available
	if t.Status != sqlc.TableStatusAvailable {
		return nil, errors.Wrap(api.ErrTableNotAvaliable.Error, "store")
	}

	return s.store.CreateDiningOrderTx(ctx, tableID, employeeID)
}

func (s *service) IsOngoingOrder(ctx context.Context, orderID pgtype.UUID) (bool, error) {
	// Check if the orderID is open and shit
	o, err := s.store.GetOrderByID(ctx, orderID)
	if err != nil {
		return false, errors.Wrap(err, "IsOngoingOrder")
	}

	return o.Status == sqlc.OrderStatusOngoing, nil
}

func (s *service) AddItemsToOrder(ctx context.Context, orderID pgtype.UUID, items []RequestItem) error {
	// Check if the orderID is open and shit
	ongoing, err := s.IsOngoingOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errors.Wrap(api.ErrUnknownOrder.Error, "store")
		}
		return errors.Wrap(err, "store")
	}

	if !ongoing {
		return errors.Wrap(api.ErrOrderNotOngoing.Error, "store")
	}

	var arg sqlc.AddOrderItemsBulkParams
	arg.OrderID = orderID

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

func (s *service) UpdateOrderItem(ctx context.Context, orderID pgtype.UUID, orderItemID int, status sqlc.OrderItemStatus) error {
	// Check if the order is valid
	ongoing, err := s.IsOngoingOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errors.Wrap(api.ErrUnknownOrder.Error, "store")
		}
		return errors.Wrap(err, "store")
	}

	if !ongoing {
		return errors.Wrap(api.ErrOrderNotOngoing.Error, "store")
	}

	// Check if the order contains the order item

	_, err = s.store.GetOrderItemByID(ctx, sqlc.GetOrderItemByIDParams{ID: int64(orderItemID), OrderID: orderID})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errors.Wrap(api.ErrUnknownOrderItem.Error, "store")
		}
		return errors.Wrap(err, "store")
	}

	arg := sqlc.UpdateOrderItemStatusParams{
		ID:      int64(orderItemID),
		OrderID: orderID,
		Status:  status,
	}

	return s.store.UpdateOrderItemStatus(ctx, arg)
}

func (s *service) GetTables(ctx context.Context, status sqlc.TableStatus) ([]Table, error) {
	t, err := s.store.GetTables(ctx, status)
	if err != nil {
		return []Table{}, errors.Wrap(err, "store")
	}

	if len(t) == 0 {
		return []Table{}, nil
	}

	var tables []Table
	for _, table := range t {
		tables = append(tables, Table{
			ID:       table.ID,
			Capacity: table.Capacity,
			Status:   table.Status,
			Notes:    table.Notes,
		})
	}

	return tables, nil
}
