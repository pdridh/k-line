package dining

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func (h *handler) CreateSession() http.HandlerFunc {

	type RequestPayload struct {
		TableID int `json:"table_id" validate:"required"`
	}

	type ResponsePayload struct {
		ID          uuid.UUID     `json:"id"`
		Status      SessionStatus `json:"status"`
		TableID     int           `json:"table_id"`
		StartedAt   time.Time     `json:"started_at"`
		CompletedAt *time.Time    `json:"completed_at,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

		// Check if table is available before creation a session on it
		available, err := h.Service.IsTableAvailable(r.Context(), p.TableID)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		if !available {
			api.WriteJSON(w, r, http.StatusConflict, "table is occupied")
			return
		}

		sess, err := h.Service.CreateSession(r.Context(), p.TableID)
		if err != nil {
			log.Println(err)
			api.WriteInternalError(w, r)
			return
		}

		// TODO change the map to a more standard way to write responses
		api.WriteJSON(w, r, http.StatusCreated, map[string]any{"session": ResponsePayload{
			ID:        sess.ID,
			Status:    sess.Status,
			TableID:   sess.TableID,
			StartedAt: sess.StartedAt,
		}})
	}

}

func (h *handler) AddItemsToSession() http.HandlerFunc {

	type RequestPayload struct {
		Items []SessionItem `json:"items" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tableIDStr := r.PathValue("tableID")
		tableID, err := strconv.Atoi(tableIDStr)
		if err != nil {
			api.WriteNotFoundError(w, r)
			return
		}

		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			log.Println(err)
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Service.Validate.Struct(p); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "Validation errors", v)
			return
		}

		i, err := h.Service.AddItemsToSession(r.Context(), tableID, p.Items)

		if err != nil {
			switch err {
			case ErrTableNoOpenSession:
				api.WriteError(w, r, http.StatusConflict, "table has no session (empty table)", nil)
				return
			default:
				log.Println(err)
				api.WriteInternalError(w, r)
				return
			}
		}

		api.WriteJSON(w, r, http.StatusCreated, i)
	}
}
