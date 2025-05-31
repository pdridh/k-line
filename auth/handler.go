package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/user"
)

type handler struct {
	Validate *validator.Validate
	Store    user.Store
}

func NewHandler(v *validator.Validate, s user.Store) *handler {
	return &handler{
		Validate: v,
		Store:    s,
	}
}

func (h *handler) RegisterUser() http.HandlerFunc {
	type RequestPayload struct {
		Name     string        `json:"name" validate:"required"`
		Type     user.UserType `json:"type" validate:"required,oneof=admin waiter kitchen"`
		Password string        `json:"password" validate:"required,min=8,max=32"`
	}

	type ResponsePayload struct {
		ID        uuid.UUID     `json:"id"`
		Name      string        `json:"name"`
		Type      user.UserType `json:"type"`
		CreatedAt time.Time     `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var p RequestPayload
		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Validate.Struct(p); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "Validation errors", v)
			return
		}

		u, err := h.Store.CreateUser(r.Context(), p.Name, p.Type, p.Password)
		if err != nil {
			log.Println(err)
			api.WriteInternalError(w, r)
			return
		}

		res := ResponsePayload{
			ID:        u.ID,
			Name:      u.Name,
			Type:      u.Type,
			CreatedAt: u.CreatedAt,
		}

		api.WriteJSON(w, r, http.StatusCreated, res)
	}
}
