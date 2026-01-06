package app

import (
	"net/http"

	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/handler"
	"github.com/gostructure/app/internal/middleware"
)

// App represents the application
type App struct {
	config  *config.Config
	router  *http.ServeMux
	handler *handler.Handler
}

// New creates a new application instance
func New(cfg *config.Config) *App {
	app := &App{
		config: cfg,
		router: http.NewServeMux(),
	}

	app.handler = handler.New(cfg)
	app.setupRoutes()

	return app
}

// Router returns the HTTP router with middleware
func (a *App) Router() http.Handler {
	// Apply middleware chain
	var h http.Handler = a.router
	h = middleware.Logging(h)
	h = middleware.Recovery(h)
	h = middleware.CORS(h)
	h = middleware.RequestID(h)

	return h
}

// setupRoutes configures all application routes
func (a *App) setupRoutes() {
	// Health check endpoints
	a.router.HandleFunc("GET /health", a.handler.Health)
	a.router.HandleFunc("GET /ready", a.handler.Ready)

	// API v1 routes
	a.router.HandleFunc("GET /api/v1/info", a.handler.Info)

	// User routes
	a.router.HandleFunc("GET /api/v1/users", a.handler.ListUsers)
	a.router.HandleFunc("GET /api/v1/users/{id}", a.handler.GetUser)
	a.router.HandleFunc("POST /api/v1/users", a.handler.CreateUser)
	a.router.HandleFunc("PUT /api/v1/users/{id}", a.handler.UpdateUser)
	a.router.HandleFunc("DELETE /api/v1/users/{id}", a.handler.DeleteUser)

	// Item routes
	a.router.HandleFunc("GET /api/v1/items", a.handler.ListItems)
	a.router.HandleFunc("GET /api/v1/items/{id}", a.handler.GetItem)
	a.router.HandleFunc("POST /api/v1/items", a.handler.CreateItem)
	a.router.HandleFunc("PUT /api/v1/items/{id}", a.handler.UpdateItem)
	a.router.HandleFunc("DELETE /api/v1/items/{id}", a.handler.DeleteItem)
}
