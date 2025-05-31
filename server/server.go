package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/pdridh/k-line/auth"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/menu"
	"github.com/pdridh/k-line/user"
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

	userStore := user.NewPSQLStore(d)
	authHandler := auth.NewHandler(v, userStore)

	menuStore := menu.NewPSQLStore(d)
	menuHandler := menu.NewHandler(v, menuStore)

	// TODO add authorization
	mux.Handle("POST /user", authHandler.HandlePostUser())

	mux.Handle("GET /menu", menuHandler.HandleGetAll())
	mux.Handle("GET /menu/{id}", menuHandler.HandleGetOne())
	mux.Handle("POST /menu", menuHandler.HandlePostMenuItem())
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
