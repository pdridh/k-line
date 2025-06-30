package dining

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db/sqlc"
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
			api.WriteValidationError(w, r, err)
			return
		}

		var userID pgtype.UUID
		if err := userID.Scan(userIDstr); err != nil {
			api.WriteInternalError(w, r)
		}

		o, err := h.Service.CreateOrder(r.Context(), p.TableID, userID)
		if err != nil {
			switch {
			case errors.Is(err, api.ErrUnknownTable.Error):
				api.WriteNotFoundError(w, r)
				return
			case errors.Is(err, api.ErrTableNotAvaliable.Error):
				api.WriteError(w, r, http.StatusConflict, api.ErrTableNotAvaliable, nil)
				return
			default:
				api.WriteInternalError(w, r)
				return
			}
		}

		api.WriteSuccess(w, r, http.StatusCreated, "Order created succefully", o)
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
			api.WriteValidationError(w, r, err)
			return
		}

		if err := h.Service.AddItemsToOrder(r.Context(), id, p.Items); err != nil {
			switch {
			case errors.Is(err, api.ErrUnknownOrder.Error):
				api.WriteNotFoundError(w, r)
				return
			case errors.Is(err, api.ErrOrderNotOngoing.Error):
				api.WriteError(w, r, http.StatusConflict, api.ErrOrderNotOngoing, nil)
				return
			default:
				api.WriteInternalError(w, r)
				return
			}
		}

		api.WriteSuccess(w, r, http.StatusCreated, "Succesfully added item to order", nil)
	}
}

func (h *handler) UpdateOrderItem() http.HandlerFunc {
	type RequestPayload struct {
		Status sqlc.OrderItemStatus `json:"status" validate:"required,oneof=pending preparing ready served delivered cancelled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := r.PathValue("order_id")
		itemIDStr := r.PathValue("item_id")

		orderID := pgtype.UUID{}
		if err := orderID.Scan(orderIDStr); err != nil {
			api.WriteNotFoundError(w, r)
			return
		}

		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			api.WriteNotFoundError(w, r)
			return
		}

		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Service.Validate.Struct(p); err != nil {
			api.WriteValidationError(w, r, err)
			return
		}

		if err := h.Service.UpdateOrderItem(r.Context(), orderID, itemID, p.Status); err != nil {
			switch {
			case errors.Is(err, api.ErrUnknownOrder.Error), errors.Is(err, api.ErrUnknownOrderItem.Error):
				api.WriteNotFoundError(w, r)
				return
			case errors.Is(err, api.ErrOrderNotOngoing.Error):
				api.WriteError(w, r, http.StatusConflict, api.ErrOrderNotOngoing, nil)
				return
			default:
				api.WriteInternalError(w, r)
				return
			}
		}

		api.WriteSuccess(w, r, http.StatusOK, "Succefully updated order item status", nil)
	}

}

func (h *handler) GetTables() http.HandlerFunc {
	type QueryParams struct {
		Status sqlc.TableStatus `json:"status" validate:"required,oneof=available occupied closed"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var p QueryParams

		api.ParseQueryParams(r.URL.Query(), &p)

		if err := h.Service.Validate.Struct(p); err != nil {
			api.WriteValidationError(w, r, err)
			return
		}

		t, err := h.Service.GetTables(r.Context(), p.Status)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteSuccess(w, r, http.StatusOK, "Retrieval successful", t)
	}
}
