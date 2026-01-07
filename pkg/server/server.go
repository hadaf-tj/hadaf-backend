package server

import (
	"context"
	"net/http"
	"shb/internal/configs"
	"shb/internal/handlers"
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

func NewServer(cfg *configs.ServerConfig, handler *handlers.Handler) *Server {
	readTimeout, _ := time.ParseDuration(cfg.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(cfg.WriteTimeout)

	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler.InitRoutes(),
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		},
	}
}
