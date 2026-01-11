package http

import (
	"database/sql"
	"net/http"

	"github.com/gostructure/app/internal/adapter/auth/bcrypt"
	"github.com/gostructure/app/internal/adapter/auth/jwt"
	"github.com/gostructure/app/internal/adapter/handler/http/user"
	"github.com/gostructure/app/internal/adapter/storage/mysql"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/core/service"
	"github.com/gostructure/app/internal/middleware"
)

func NewServer(cfg *config.Config, db *sql.DB) (*Server, error) {
	// repositories
	userRepo := mysql.NewUserRepository(db)
	tokenRepo := mysql.NewTokenRepository(db)

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
