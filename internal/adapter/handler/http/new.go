package http

import (
	"fmt"
	"net/http"

	"github.com/gostructure/app/internal/adapter/auth/bcrypt"
	"github.com/gostructure/app/internal/adapter/auth/jwt"
	"github.com/gostructure/app/internal/adapter/handler/http/user"
	"github.com/gostructure/app/internal/adapter/storage"
	"github.com/gostructure/app/internal/adapter/storage/mysql"
	"github.com/gostructure/app/internal/adapter/storage/postgres"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/core/port"
	"github.com/gostructure/app/internal/core/service"
	"github.com/gostructure/app/internal/middleware"
)

func NewServer(cfg *config.Config, db storage.Database) (*Server, error) {
	var userRepo port.UserRepository
	var tokenRepo port.TokenRepository

	sqlDB := db.GetDB()

	// 1. Initialize repositories based on driver
	switch db.DriverName() {
	case "mysql":
		userRepo = mysql.NewUserRepository(sqlDB)
		tokenRepo = mysql.NewTokenRepository(sqlDB)
	case "postgres":
		userRepo = postgres.NewUserRepository(sqlDB)
		tokenRepo = postgres.NewTokenRepository(sqlDB)
	default:
		return nil, fmt.Errorf("unsupported database driver in server: %s", db.DriverName())
	}

	// auth
	passwordHasher := bcrypt.NewBcryptHasher()
	jwtProvider := jwt.NewJWTProvider(&cfg.JWT)

	// services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		passwordHasher,
		jwtProvider,
	)

	// handlers
	userHandler := user.NewUserHandler(userService)
	authHandler := user.NewAuthHandler(authService)

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtProvider, userRepo)

	handlers := &Handlers{
		User: userHandler,
		Auth: authHandler,
	}

	mux := http.NewServeMux()

	server := &Server{
		cfg:            cfg,
		db:             db,
		mux:            mux,
		handlers:       handlers,
		authMiddleware: authMiddleware,
	}

	server.registerRoutes(mux)

	return server, nil
}
