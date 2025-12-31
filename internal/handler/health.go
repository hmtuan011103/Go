package handler

import (
	"net/http"

	"github.com/gostructure/app/pkg/response"
)

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// Ready handles readiness check requests
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	// Add database connectivity check here if needed
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

// Info returns application information
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"name":        h.config.App.Name,
		"version":     h.config.App.Version,
		"environment": h.config.App.Environment,
	})
}
