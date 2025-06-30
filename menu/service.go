package menu

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/db/sqlc"
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

func (s *service) CreateItem(ctx context.Context, name string, description pgtype.Text, price float64, requiresTicket bool) (*Item, error) {

	arg := sqlc.CreateMenuItemParams{
		Name:           name,
		Description:    description,
		Price:          price,
		RequiresTicket: requiresTicket,
	}

	i, err := s.store.CreateMenuItem(ctx, arg)
	if err != nil {
		errCode := db.GetSQLErrorCode(err)
		if errCode == db.UniqueViolation {
			return nil, api.ErrItemNameConflict.Error
		}
		return nil, err
	}

	item := &Item{
		ID:             i.ID,
		Price:          i.Price,
		Name:           i.Name,
		Description:    i.Description,
		RequiresTicket: i.RequiresTicket,
		CreatedAt:      i.CreatedAt,
	}

	return item, nil
}

func (s *service) GetItems(ctx context.Context, search string, limit int32, offset int32) ([]Item, error) {
	arg := sqlc.GetMenuItemsParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	}

	i, err := s.store.GetMenuItems(ctx, arg)
	if err != nil {
		return []Item{}, err
	}

	var items []Item
	for _, item := range i {
		items = append(items, Item{
			ID:             item.ID,
			Name:           item.Name,
			Description:    item.Description,
			Price:          item.Price,
			CreatedAt:      item.CreatedAt,
			RequiresTicket: item.RequiresTicket,
		})
	}

	return items, nil
}

func (s *service) GetItemByID(ctx context.Context, id int32) (*Item, error) {

	i, err := s.store.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, api.ErrUnkownMenuItem.Error
		}

		return nil, err
	}

	item := &Item{
		ID:             i.ID,
		Name:           i.Name,
		Description:    i.Description,
		Price:          i.Price,
		CreatedAt:      i.CreatedAt,
		RequiresTicket: i.RequiresTicket,
	}

	return item, nil
}
