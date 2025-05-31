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

	authService := auth.NewService(v, userStore)
	authHandler := auth.NewHandler(authService)

	menuStore := menu.NewPSQLStore(d)
	menuHandler := menu.NewHandler(v, menuStore)

	// TODO add authorization for different user types

	mux.Handle("POST /auth/register", authHandler.Register())
	mux.Handle("POST /auth/login", authHandler.Login())

	mux.Handle("GET /menu", menuHandler.GetAllItems())
	mux.Handle("GET /menu/{id}", menuHandler.GetItemById())
	mux.Handle("POST /menu", menuHandler.CreateItem())
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
