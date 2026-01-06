package main

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gostructure/app/internal/adapter/handler/http"
	"github.com/gostructure/app/internal/adapter/storage/memory"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/middleware"
	"github.com/gostructure/app/internal/service"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Starting %s on %s (%s)\n", cfg.App.Name, cfg.Server.Address, cfg.App.Environment)

	// 2. Setup Adapters (Repository)
	userRepo := memory.NewUserRepository()

	// 3. Setup Core (Service)
	userService := service.NewUserService(userRepo)

	// 4. Setup Adapters (Handler)
	userHandler := http.NewUserHandler(userService)

	// 5. Setup Router (Standard Mux)
	router := stdMux(userHandler)

	// Middleware
	handler := middleware.Logging(router)

	// 6. Setup HTTP Server
	serverAddr := cfg.Server.Address
	srv := &nethttp.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 7. Start Server in a goroutine
	go func() {
		log.Printf("Server listening on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until signal is received

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// stdMux setups the standard library router
// In a real app, this might be in a separate router.go file
func stdMux(userHandler *http.UserHandler) *nethttp.ServeMux {
	mux := nethttp.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("OK"))
	})

	// User Routes (Go 1.22+ patterns)
	mux.HandleFunc("GET /api/v1/users", userHandler.ListUsers)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /api/v1/users", userHandler.CreateUser)
	mux.HandleFunc("PUT /api/v1/users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/v1/users/{id}", userHandler.DeleteUser)

	return mux
}

// Alias for net/http to avoid conflict with package name "http"
// inside main since we imported internal/adapter/handler/http
// We can rename import instead.
