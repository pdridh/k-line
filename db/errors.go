package db

import "github.com/pkg/errors"

var (
	ErrAlreadyConnected = errors.New("already connected")
	ErrNotConnected     = errors.New("database is not connected")
)
