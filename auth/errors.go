package auth

import "github.com/pkg/errors"

var (
	ErrUnknownEmail  = errors.New("unkown email")
	ErrWrongPassword = errors.New("wrong password")
)
