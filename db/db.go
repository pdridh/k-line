package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/config"
)

var core *sqlx.DB

// Connect to the db and store it as core
func Connect() (*sqlx.DB, error) {

	if core != nil {
		return nil, ErrAlreadyConnected
	}

	driver := "postgres"
	uri := config.Server().DatabaseURI

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

func Disconnect() error {
	if core == nil {
		return nil
	}

	log.Println("Closing database connection")
	return core.Close()
}
