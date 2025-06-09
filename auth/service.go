package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/db/sqlc"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	Validate *validator.Validate
	Store    db.Store
}

// Simple wrapper to create a new user service given the stores and validator
func NewService(v *validator.Validate, u db.Store) *service {
	return &service{
		Validate: v,
		Store:    u,
	}
}

func (s *service) CreateUser(ctx context.Context, email string, name string, userType sqlc.UserType, password string) (sqlc.User, error) {
	// Hash password
	var u sqlc.User
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return u, errors.Wrap(err, "hash")
	}

	arg := sqlc.CreateUserParams{
		Email:    email,
		Name:     name,
		Type:     userType,
		Password: hashedPassword,
	}

	u, err = s.Store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.GetSQLErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			return u, errors.Wrap(ErrEmailAlreadyExists, "store")
		}
	}

	return u, nil
}

func (s *service) AuthenticateUser(ctx context.Context, email string, password string) (string, *sqlc.User, error) {
	u, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return "", nil, errors.Wrap(ErrUnknownEmail, "store")
		}
		return "", nil, errors.Wrap(err, "store")
	}

	// Check if the password is correct
	if err := CompareHashedPasswords(u.Password, password); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return "", nil, ErrWrongPassword
		default:
			return "", nil, err
		}
	}

	t, err := GenerateJWT(u.ID.String(), u.Email, u.Name, u.Type, config.Server().JWTExpiration)
	if err != nil {
		return "", nil, errors.Wrap(err, "jwtgen")
	}

	return t, &u, nil
}
