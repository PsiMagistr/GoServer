package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}
type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *HTTPServer) Run() error {
	fmt.Printf("Сервер запускается на %s... \n", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	fmt.Println("Остановка сервера...")
	return s.server.Shutdown(ctx)
}
