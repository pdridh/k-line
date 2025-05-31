package user

import "github.com/pkg/errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)
