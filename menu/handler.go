package menu

import (
	"net/http"

	"github.com/pdridh/k-line/api"
)

type handler struct {
	Store Store
}

func NewHandler(s Store) *handler {
	return &handler{
		Store: s,
	}
}

func (h *handler) HandlePostMenuItem() http.HandlerFunc {
	type RequestPayload struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload RequestPayload

		if err := api.ParseJSON(r, &payload); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		// TODO validate

		i, err := h.Store.CreateItem(r.Context(), payload.Name, payload.Description, payload.Price)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusCreated, i)
	}

}
