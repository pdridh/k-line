package auth

import "github.com/pkg/errors"

var (
	ErrUnknownEmail               = errors.New("unknown email")
	ErrWrongPassword              = errors.New("wrong password")
	ErrUnexpectedJWTSigningMethod = errors.New("unexpected signing method")
	ErrEmailAlreadyExists         = errors.New("email conflict in create")
)
