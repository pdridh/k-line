package dining

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
)

type handler struct {
	Service *service
}

func NewHandler(s *service) *handler {
	return &handler{
		Service: s,
	}
}

func (h *handler) CreateOrder() http.HandlerFunc {

	type RequestPayload struct {
		TableID pgtype.Text `json:"table_id" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userIDstr := api.CurrentUserID(r)

		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Service.Validate.Struct(p); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "Validation errors", v)
			return
		}

		var userID pgtype.UUID
		if err := userID.Scan(userIDstr); err != nil {
			api.WriteInternalError(w, r)
		}

		o, err := h.Service.CreateOrder(r.Context(), p.TableID, userID)
		if err != nil {
			switch {
			case errors.Is(err, ErrUnknownTable):
				api.WriteNotFoundError(w, r)
				return
			case errors.Is(err, ErrTableNotAvaliable):
				api.WriteError(w, r, http.StatusConflict, "table is not available", nil)
				return
			default:
				log.Println(err)
				api.WriteInternalError(w, r)
				return
			}
		}

		api.WriteJSON(w, r, http.StatusCreated, map[string]any{"order": o})
	}

}

func (h *handler) AddOrderItem() http.HandlerFunc {

	type RequestPayload struct {
		Items []RequestItem `json:"items" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		// Check if order is valid
		id := pgtype.UUID{}
		if err := id.Scan(idStr); err != nil {
			api.WriteNotFoundError(w, r)
			return
		}

		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Service.Validate.Struct(p); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "Validation errors", v)
			return
		}

		h.Service.AddItemsToOrder(r.Context(), id, p.Items)

		api.WriteJSON(w, r, http.StatusCreated, p)
	}
}
