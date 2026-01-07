package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gostructure/app/internal/core/port"
	"github.com/gostructure/app/pkg/response"
)

type AuthMiddleware struct {
	tokenProv port.TokenProvider
	userRepo  port.UserRepository
}

func NewAuthMiddleware(tokenProv port.TokenProvider, userRepo port.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenProv: tokenProv,
		userRepo:  userRepo,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := m.tokenProv.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// CHECK TOKEN VERSION (Security Stamp)
		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil || user.TokenVersion != claims.Version {
			response.Error(w, http.StatusUnauthorized, "Token has been invalidated (logged in elsewhere or expired)")
			return
		}

		// Inject user info into context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
