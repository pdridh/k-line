package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdridh/k-line/db/sqlc"
)

type UserClaims struct {
	UserID    string
	UserEmail string
	UserName  string
	UserType  sqlc.UserType
	jwt.RegisteredClaims
}

type User struct {
	ID        pgtype.UUID      `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	Type      sqlc.UserType    `json:"type"`
	CreatedAt pgtype.Timestamp `json:"created_at,omitempty"`
}

type UserAuth struct {
	ID    string        `json:"id"`
	Email string        `json:"email"`
	Name  string        `json:"name"`
	Type  sqlc.UserType `json:"type"`
}
