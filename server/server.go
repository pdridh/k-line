package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/menu"
)

const (
	WriteTimeout = 15 * time.Second
	ReadTimeout  = 15 * time.Second
)

type server struct {
	HttpServer *http.Server
}

func New(v *validator.Validate, d *sqlx.DB) *server {
	mux := http.NewServeMux()

	menuStore := menu.NewPSQLStore(d)
	menuHandler := menu.NewHandler(v, menuStore)

	mux.Handle("GET /menu/items", menuHandler.HandleGetAll())
	mux.Handle("POST /menu/items", menuHandler.HandlePostMenuItem())
	mux.Handle("/", http.NotFoundHandler())

	h := &http.Server{
		Addr:         net.JoinHostPort(config.Server().Host, config.Server().Port),
		Handler:      mux,
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
	}

	return &server{HttpServer: h}
}

func (s *server) Start() error {
	log.Println("Listening on: ", s.HttpServer.Addr)

	if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
