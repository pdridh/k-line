package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pdridh/k-line/db/sqlc"
)

type UserClaims struct {
	UserID   string
	UserType sqlc.UserType
	jwt.RegisteredClaims
}
