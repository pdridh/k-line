package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pdridh/k-line/db/sqlc"
)

type Store interface {
	sqlc.Querier
	CreateDiningOrderTx(ctx context.Context, tableID pgtype.Text, employeeID pgtype.UUID) (*pgtype.UUID, error)
}

type psqlStore struct {
	pool *pgxpool.Pool
	*sqlc.Queries
}

func NewPSQLStore(pool *pgxpool.Pool) Store {
	return &psqlStore{
		pool:    pool,
		Queries: sqlc.New(pool),
	}
}

func (s *psqlStore) CreateDiningOrderTx(ctx context.Context, tableID pgtype.Text, employeeID pgtype.UUID) (*pgtype.UUID, error) {

	var orderID pgtype.UUID
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		orderID, err = q.CreateOrder(ctx, sqlc.CreateOrderParams{
			Type:       sqlc.OrderTypeDining,
			TableID:    tableID,
			EmployeeID: employeeID,
		})

		if err != nil {
			return err
		}

		updateArg := sqlc.UpdateTableStatusParams{
			Status: sqlc.TableStatusOccupied,
			ID:     tableID.String,
		}

		if err := q.UpdateTableStatus(ctx, updateArg); err != nil {
			return err
		}

		return nil
	})

	return &orderID, err
}
