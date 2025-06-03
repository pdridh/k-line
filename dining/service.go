package dining

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type service struct {
	Validate *validator.Validate
	store    Store
}

func NewService(v *validator.Validate, s Store) *service {
	return &service{
		Validate: v,
		store:    s,
	}
}

func (s *service) CreateSession(ctx context.Context, tableID int) (*Session, error) {
	return s.store.CreateSession(ctx, tableID)
}

func (s *service) IsTableAvailable(ctx context.Context, tableID int) (bool, error) {

	sess, err := s.store.GetOngoingSessionByTable(ctx, tableID)
	if err != nil {
		return false, errors.Wrap(err, "store")
	}

	return sess == nil, nil
}

func (s *service) GetOngoingTable(ctx context.Context, tableID int) (*Session, error) {
	return s.store.GetOngoingSessionByTable(ctx, tableID)
}

func (s *service) AddItemsToSession(ctx context.Context, tableID int, items []SessionItem) ([]SessionItem, error) {
	sess, err := s.GetOngoingTable(ctx, tableID)

	if err != nil {
		return nil, errors.Wrap(err, "store")
	}

	if sess == nil {
		return nil, ErrTableNoOpenSession
	}

	return s.store.CreateSessionItems(ctx, sess.ID, items)
}
