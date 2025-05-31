package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	Validate  *validator.Validate
	UserStore user.Store
}

// Simple wrapper to create a new user service given the stores and validator
func NewService(v *validator.Validate, u user.Store) *service {
	return &service{
		Validate:  v,
		UserStore: u,
	}
}

func (s *service) CreateUser(ctx context.Context, email string, name string, userType user.UserType, password string) (*user.User, error) {
	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "hash")
	}

	u, err := s.UserStore.Create(ctx, email, name, userType, hashedPassword)
	if err != nil {
		return nil, errors.Wrap(err, "store")
	}

	return u, nil
}

func (s *service) AuthenticateUser(ctx context.Context, email string, password string) (string, error) {
	u, err := s.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrap(err, "store")
	}

	if u == nil {
		return "", ErrUnknownEmail
	}

	// Check if the password is correct
	if err := CompareHashedPasswords(u.Password, password); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return "", ErrWrongPassword
		default:
			return "", err
		}
	}

	t, err := GenerateJWT(u.ID.String(), u.Type, config.Server().JWTExpiration)
	if err != nil {
		return "", errors.Wrap(err, "jwtgen")
	}

	return t, nil
}
