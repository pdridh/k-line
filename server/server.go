package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pdridh/k-line/config"
)

const (
	WriteTimeout = 15 * time.Second
	ReadTimeout  = 15 * time.Second
)

type server struct {
	HttpServer *http.Server
}

func New() *server {
	mux := http.NewServeMux()

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
	if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	log.Println("Listening on: ", s.HttpServer.Addr)

	return nil
}
