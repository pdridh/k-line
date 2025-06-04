package auth

import (
	"context"
	"net/http"
	"slices"

	"github.com/pdridh/k-line/api"
	"github.com/pdridh/k-line/db/sqlc"
)

// Takes a handler function and only calls it if
// the jwt token it extracts from the request's is valid.
// The next handler function is called with the userid in context
func Middleware(next http.HandlerFunc, allowedTypes ...sqlc.UserType) http.HandlerFunc {
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

		newCtx := context.WithValue(r.Context(), api.ContextUserKey, api.CurrentUser{ID: c.UserID, Type: c.UserType})
		if c.UserType == sqlc.UserTypeAdmin || slices.Contains(allowedTypes, c.UserType) {
			next(w, r.WithContext(newCtx))
			return
		}

		api.WriteError(w, r, http.StatusForbidden, "not allowed to use this route", nil)
	}
}
