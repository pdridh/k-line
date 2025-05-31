package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pdridh/k-line/user"
)

type UserClaims struct {
	UserID   string
	UserType user.UserType
	jwt.RegisteredClaims
}
