package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/server"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config.Load()

	uri := config.Server().DatabaseURI

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		log.Fatalln("cannot create connection pool", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalln("cannot connect to db:", err)
	}

	store := db.NewPSQLStore(pool)

	v := validator.New()
	s := server.New(v, store)

	if err := s.Start(); err != nil {
		log.Fatalln("failed to start the server: ", err)
	}
}
