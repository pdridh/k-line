package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pdridh/k-line/db/sqlc"
)

type UserClaims struct {
	UserID    string
	UserEmail string
	UserName  string
	UserType  sqlc.UserType
	jwt.RegisteredClaims
}
