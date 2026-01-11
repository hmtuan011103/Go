package http

import (
	"database/sql"
	"net/http"

	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/middleware"
)

type Server struct {
	cfg            *config.Config
	db             *sql.DB
	mux            *http.ServeMux
	handlers       *Handlers
	authMiddleware *middleware.AuthMiddleware
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
