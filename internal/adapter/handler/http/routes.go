package http

import nethttp "net/http"

func (s *Server) registerRoutes(mux *nethttp.ServeMux) {
	// Health
	mux.HandleFunc("GET /health", func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth
	mux.HandleFunc("POST /api/v1/auth/register", s.handlers.Auth.Register)
	mux.HandleFunc("POST /api/v1/auth/login", s.handlers.Auth.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", s.handlers.Auth.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", s.handlers.Auth.Logout)

	// User (Protected)
	protected := s.authMiddleware.Handle

	mux.Handle(
		"GET /api/v1/users",
		protected(nethttp.HandlerFunc(s.handlers.User.ListUsers)),
	)

	mux.Handle(
		"GET /api/v1/users/{id}",
		protected(nethttp.HandlerFunc(s.handlers.User.GetUser)),
	)

	mux.Handle(
		"POST /api/v1/users",
		protected(nethttp.HandlerFunc(s.handlers.User.CreateUser)),
	)

	mux.Handle(
		"PUT /api/v1/users/{id}",
		protected(nethttp.HandlerFunc(s.handlers.User.UpdateUser)),
	)

	mux.Handle(
		"DELETE /api/v1/users/{id}",
		protected(nethttp.HandlerFunc(s.handlers.User.DeleteUser)),
	)
}
