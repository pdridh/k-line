package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db"
)

func main() {
	config.Load()
	_, err := db.Connect()
	if err != nil {
		log.Println("failed to connect to db: ", err)
		db.Disconnect()
		return
	}
}
