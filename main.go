package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/server"
)

func main() {
	config.Load()

	d, err := db.Connect()
	if err != nil {
		log.Println("failed to connect to db: ", err)
		return
	}

	v := validator.New()
	s := server.New(v, d)

	if err := s.Start(); err != nil {
		log.Fatalln("Failed to start the server: ", err)
	}
}
