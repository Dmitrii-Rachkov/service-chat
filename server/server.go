package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	"service-chat/internal/config"
)

// Server - структура сервера из пакета http
type Server struct {
	httpServer *http.Server
}

// Run - запуск сервера
func (s *Server) Run(cfg *config.Config, router *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return s.httpServer.ListenAndServe()
}

// ShutDown - остановка сервера
func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
