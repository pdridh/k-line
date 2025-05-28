package main

import (
	"log"

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

	s := server.New(d)
	if err := s.Start(); err != nil {
		log.Fatalln("Failed to start the server: ", err)
	}
}
