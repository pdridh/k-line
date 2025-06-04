package api

import (
	"net/http"

	"github.com/pdridh/k-line/db/sqlc"
)

type contextKey string

const ContextUserKey contextKey = "user"

type CurrentUser struct {
	ID   string
	Type sqlc.UserType
}

// Given a request extracts the value of the userID (string) from the context using the ContextUserKey
func CurrentUserID(r *http.Request) string {
	return r.Context().Value(ContextUserKey).(CurrentUser).ID
}

// Given a request extracts the value of the CurrentUser.Type (string) from the context using the ContextUserKey
func CurrentUserType(r *http.Request) sqlc.UserType {
	return r.Context().Value(ContextUserKey).(CurrentUser).Type
}
