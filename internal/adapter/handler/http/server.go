package http

import (
	"net/http"

	"github.com/gostructure/app/internal/adapter/storage"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/middleware"
)

type Server struct {
	cfg            *config.Config
	db             storage.Database
	mux            *http.ServeMux
	handlers       *Handlers
	authMiddleware *middleware.AuthMiddleware
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
