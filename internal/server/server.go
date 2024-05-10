package server

import (
	"context"
	"fmt"
	"net/http"
)

type Option struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type Server struct {
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(handler http.Handler, opt Option) *Server {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Handler: handler,
	}
	return &Server{
		handler:    handler,
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}
