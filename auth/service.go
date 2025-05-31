package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
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

func (s *service) AuthenticateUser(ctx context.Context, email string, password string) (bool, error) {
	u, err := s.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "store")
	}

	if u == nil {
		return false, ErrUnknownEmail
	}

	// Check if the password is correct
	if err := CompareHashedPasswords(u.Password, password); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, ErrWrongPassword
		default:
			return false, err
		}
	}

	return true, nil
}
