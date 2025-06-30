package menu

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/db/sqlc"
)

type handler struct {
	Store    db.Store
	Validate *validator.Validate
}

func NewHandler(v *validator.Validate, s db.Store) *handler {
	return &handler{
		Validate: v,
		Store:    s,
	}
}

func (h *handler) CreateItem() http.HandlerFunc {
	type RequestPayload struct {
		Name           string      `json:"name" validate:"required"`
		Description    pgtype.Text `json:"description" validate:"required"`
		Price          float64     `json:"price" validate:"required"`
		RequiresTicket bool        `json:"requires_ticket" validate:"required"`
	}

	type ResponsePayload struct {
		ID          int32            `json:"id"`
		Name        string           `json:"name"`
		Description pgtype.Text      `json:"description"`
		Price       float64          `json:"price"`
		CreatedAt   pgtype.Timestamp `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload RequestPayload

		if err := api.ParseJSON(r, &payload); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Validate.Struct(payload); err != nil {
			api.WriteValidationError(w, r, err)
			return
		}

		arg := sqlc.CreateMenuItemParams{
			Name:           payload.Name,
			Description:    payload.Description,
			Price:          payload.Price,
			RequiresTicket: payload.RequiresTicket,
		}
		i, err := h.Store.CreateMenuItem(r.Context(), arg)
		if err != nil {
			errCode := db.GetSQLErrorCode(err)
			if errCode == db.UniqueViolation {
				api.WriteError(w, r, http.StatusConflict, api.ErrItemNameConflict, nil)
				return
			}
			api.WriteInternalError(w, r)
			return
		}

		res := ResponsePayload{
			ID:          i.ID,
			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
			CreatedAt:   i.CreatedAt,
		}

		api.WriteSuccess(w, r, http.StatusCreated, "Created new menu item", res)
	}

}

func (h *handler) GetAllItems() http.HandlerFunc {

	type ResponsePayload struct {
		ID          int32            `json:"id"`
		Name        string           `json:"name"`
		Description pgtype.Text      `json:"description"`
		Price       float64          `json:"price"`
		CreatedAt   pgtype.Timestamp `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var filters MenuFilters

		api.ParseQueryParams(r.URL.Query(), &filters)
		filters.Validate(50, 20)

		offset := (filters.Page - 1) * filters.Limit
		arg := sqlc.GetMenuItemsParams{
			Search: filters.Search,
			Limit:  filters.Limit,
			Offset: offset,
		}

		i, err := h.Store.GetMenuItems(r.Context(), arg)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		var itemsRes []ResponsePayload

		for _, item := range i {
			itemsRes = append(itemsRes, ResponsePayload{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
				CreatedAt:   item.CreatedAt,
			})
		}

		api.WriteSuccess(w, r, http.StatusOK, "Retrieval successful", itemsRes)
	}
}

func (h *handler) GetItemById() http.HandlerFunc {

	type ResponsePayload struct {
		ID          int32            `json:"id"`
		Name        string           `json:"name"`
		Description pgtype.Text      `json:"description"`
		Price       float64          `json:"price"`
		CreatedAt   pgtype.Timestamp `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			api.WriteNotFoundError(w, r)
			return
		}

		i, err := h.Store.GetItemByID(r.Context(), int32(id))
		if err != nil {

			if errors.Is(err, db.ErrRecordNotFound) {
				api.WriteNotFoundError(w, r)
				return
			}

			api.WriteInternalError(w, r)
			return
		}

		res := ResponsePayload{
			ID:          i.ID,
			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
			CreatedAt:   i.CreatedAt,
		}

		api.WriteSuccess(w, r, http.StatusOK, "Retrieval successful", res)
	}
}
