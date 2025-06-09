package auth

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db/sqlc"
	"github.com/pkg/errors"
)

type handler struct {
	Service *service
}

func NewHandler(s *service) *handler {
	return &handler{
		Service: s,
	}
}

func (h *handler) Register() http.HandlerFunc {
	type RequestPayload struct {
		Name     string        `json:"name" validate:"required"`
		Email    string        `json:"email" validate:"required,email"`
		Type     sqlc.UserType `json:"type" validate:"required,oneof=admin waiter kitchen"`
		Password string        `json:"password" validate:"required,min=8,max=32"`
	}

	type ResponsePayload struct {
		ID        pgtype.UUID      `json:"id"`
		Email     string           `json:"email"`
		Name      string           `json:"name"`
		Type      sqlc.UserType    `json:"type"`
		CreatedAt pgtype.Timestamp `json:"created_at"`
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

		u, err := h.Service.CreateUser(r.Context(), p.Email, p.Name, p.Type, p.Password)
		if errors.Is(err, ErrEmailAlreadyExists) {
			api.WriteError(w, r, http.StatusConflict, "cannot use this email", nil)
			return
		}

		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		res := ResponsePayload{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Type:      u.Type,
			CreatedAt: u.CreatedAt,
		}

		api.WriteJSON(w, r, http.StatusCreated, res)
	}
}

func (h *handler) Login() http.HandlerFunc {
	type RequestPayload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	type UserResponse struct {
		ID    pgtype.UUID   `json:"id"`
		Email string        `json:"email"`
		Name  string        `json:"name"`
		Type  sqlc.UserType `json:"type"`
	}

	type ResponsePayload struct {
		Message string       `json:"message"`
		User    UserResponse `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteBadRequestError(w, r)
			return
		}

		if err := h.Service.Validate.Struct(p); err != nil {
			v := api.FormatValidationErrors(err)
			api.WriteError(w, r, http.StatusBadRequest, "validation errors", v)
			return
		}

		t, u, err := h.Service.AuthenticateUser(r.Context(), p.Email, p.Password)
		if err != nil {
			switch {
			case errors.Is(err, ErrUnknownEmail), errors.Is(err, ErrWrongPassword):
				api.WriteError(w, r, http.StatusUnauthorized, "invalid credentials", nil)
				return
			default:
				api.WriteInternalError(w, r)
				return
			}
		}

		SetJWTCookie(w, t)

		api.WriteJSON(w, r, http.StatusOK, ResponsePayload{Message: "Login succesfull!", User: UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Type:  u.Type,
		}})
	}
}

func (h *handler) GetAuth() http.HandlerFunc {
	type ResponsePayload struct {
		UserID    string        `json:"id"`
		UserEmail string        `json:"email"`
		UserName  string        `json:"string"`
		UserType  sqlc.UserType `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		jCookie, err := r.Cookie("jwt")
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		j := jCookie.Value

		t, err := ValidateJWT(j)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		c, err := UserClaimsFromJWT(t)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, ResponsePayload{
			UserID:    c.UserID,
			UserEmail: c.UserEmail,
			UserName:  c.UserName,
			UserType:  c.UserType,
		})
	}
}
