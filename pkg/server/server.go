package server

import (
	"context"
	"net/http"
	"shb/internal/handlers"
	"shb/pkg/configs"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(cfg *configs.Server, handler *handlers.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler.InitRoutes(),
			ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		},
	}
}
