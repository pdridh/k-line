package user

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserAdmin   UserType = "admin"
	UserWaiter  UserType = "waiter"
	UserKitchen UserType = "kitchen"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Type      UserType  `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
