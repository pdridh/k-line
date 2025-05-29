package menu

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/k-line/api"
)

type handler struct {
	Store    Store
	Validate *validator.Validate
}

func NewHandler(v *validator.Validate, s Store) *handler {
	return &handler{
		Validate: v,
		Store:    s,
	}
}

func (h *handler) HandlePostMenuItem() http.HandlerFunc {
	type RequestPayload struct {
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description" validate:"required"`
		Price       float64 `json:"price" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload RequestPayload

		if err := api.ParseJSON(r, &payload); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Validate.Struct(payload); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "Validation errors", v)
			return
		}

		i, err := h.Store.CreateItem(r.Context(), payload.Name, payload.Description, payload.Price)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusCreated, i)
	}

}

func (h *handler) HandleGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		i, err := h.Store.GetAllItems(r.Context())
		if err != nil {
			log.Println(err)
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, i)
	}
}
