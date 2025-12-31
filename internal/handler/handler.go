package handler

import (
	"github.com/gostructure/app/internal/config"
)

// Handler contains all HTTP handlers
type Handler struct {
	config *config.Config
}

// New creates a new Handler instance
func New(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
	}
}
