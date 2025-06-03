package dining

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/k-line/db"
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
