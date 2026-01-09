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

	"github.com/gostructure/app/internal/adapter/auth/bcrypt"
	"github.com/gostructure/app/internal/adapter/auth/jwt"
	"github.com/gostructure/app/internal/adapter/handler/http/user"
	"github.com/gostructure/app/internal/adapter/storage/mysql"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/core/service"
	"github.com/gostructure/app/internal/middleware"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// 2. Set Timezone dynamic from config
	loc, err := time.LoadLocation(cfg.App.Timezone)
	if err == nil {
		time.Local = loc
	} else {
		log.Printf("Warning: Failed to load timezone %s, using local: %v", cfg.App.Timezone, err)
	}

	fmt.Printf("Starting %s on %s (%s)\n", cfg.App.Name, cfg.Server.Address, cfg.App.Environment)

	// 3. Setup Adapters (Repository & Auth)
	db, err := mysql.NewMySQLConnection(&cfg.Database, cfg.App.Timezone)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run Migrations
	if err := mysql.RunMigrations(db, &cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := mysql.NewUserRepository(db)
	tokenRepo := mysql.NewTokenRepository(db)
	passwordHasher := bcrypt.NewBcryptHasher()
	jwtProvider := jwt.NewJWTProvider(&cfg.JWT)

	// 3. Setup Core (Service)
	// UserService
	userService := service.NewUserService(userRepo)
	// AuthService
	authService := service.NewAuthService(userRepo, tokenRepo, passwordHasher, jwtProvider)

	// 4. Setup Adapters (Handler)
	userHandler := user.NewUserHandler(userService)
	authHandler := user.NewAuthHandler(authService)

	// Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtProvider, userRepo)

	// 5. Setup Router (Standard Mux)
	router := stdMux(userHandler, authHandler, authMiddleware)

	// Global Middleware
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
func stdMux(
	userHandler *user.UserHandler,
	authHandler *user.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *nethttp.ServeMux {
	mux := nethttp.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth Routes
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// User Routes (Protected)
	// We wrap the sensitive handlers with AuthMiddleware
	// For Go 1.22 mux, we can wrap individual handlers
	// Note: HandleFunc takes a function, Handle takes a Handler interface.
	// Middleware returns Handler.

	// Helper to wrap
	protect := authMiddleware.Handle

	mux.Handle("GET /api/v1/users", protect(nethttp.HandlerFunc(userHandler.ListUsers)))
	mux.Handle("GET /api/v1/users/{id}", protect(nethttp.HandlerFunc(userHandler.GetUser)))
	// POST users might be admin only or open? Usually creating a user is "Register" now.
	// But let's keep it protected as an "Admin create user" feature for now, or just protected.
	mux.Handle("POST /api/v1/users", protect(nethttp.HandlerFunc(userHandler.CreateUser)))
	mux.Handle("PUT /api/v1/users/{id}", protect(nethttp.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("DELETE /api/v1/users/{id}", protect(nethttp.HandlerFunc(userHandler.DeleteUser)))

	return mux
}
