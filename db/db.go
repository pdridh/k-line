package db

import (
	"os"

	"github.com/jmoiron/sqlx"
)

var core *sqlx.DB

// Connect to the db and store it as core
func Connect() (*sqlx.DB, error) {

	if core != nil {
		return nil, ErrAlreadyConnected
	}

	driver := "postgres"
	uri := os.Getenv("DB_URI")

	d, err := sqlx.Connect(driver, uri)
	if err != nil {
		return nil, err
	}

	core = d
	return core, nil
}

// Wrapper around the core db instance that protects against trying to use core before connecting
func DB() (*sqlx.DB, error) {
	if core == nil {
		return nil, ErrNotConnected
	}

	return core, nil
}
