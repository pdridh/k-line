package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/k-line/auth"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db"
	"github.com/pdridh/k-line/db/sqlc"
	"github.com/pdridh/k-line/dining"
	"github.com/pdridh/k-line/menu"
)

const (
	WriteTimeout = 15 * time.Second
	ReadTimeout  = 15 * time.Second
)

type server struct {
	HttpServer *http.Server
}

func New(v *validator.Validate, store db.Store) *server {
	mux := http.NewServeMux()

	authService := auth.NewService(v, store)
	authHandler := auth.NewHandler(authService)

	menuHandler := menu.NewHandler(v, store)

	diningService := dining.NewService(v, store)
	diningHandler := dining.NewHandler(diningService)

	mux.Handle("POST /auth/register", authHandler.Register())
	mux.Handle("POST /auth/login", authHandler.Login())

	mux.Handle("GET /menu", auth.Middleware(menuHandler.GetAllItems(), sqlc.UserTypeWaiter, sqlc.UserTypeKitchen))
	mux.Handle("GET /menu/{id}", auth.Middleware(menuHandler.GetItemById(), sqlc.UserTypeWaiter, sqlc.UserTypeKitchen))
	mux.Handle("POST /menu", auth.Middleware(menuHandler.CreateItem()))

	mux.Handle("POST /dining", auth.Middleware(diningHandler.CreateOrder(), sqlc.UserTypeWaiter))
	mux.Handle("POST /dining/{id}/item", auth.Middleware(diningHandler.AddOrderItem(), sqlc.UserTypeWaiter))
	mux.Handle("PATCH /dining/{order_id}/{item_id}", auth.Middleware(diningHandler.UpdateOrderItem(), sqlc.UserTypeWaiter, sqlc.UserTypeKitchen))

	mux.Handle("/", http.NotFoundHandler())

	var handler http.Handler = mux

	h := &http.Server{
		Addr:         net.JoinHostPort(config.Server().Host, config.Server().Port),
		Handler:      handler,
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
