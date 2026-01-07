package port

import (
	"context"

	"github.com/gostructure/app/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthResult struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) error
}

type TokenProvider interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(token string) (*Claims, error)
}

type Claims struct {
	UserID  int64
	Role    string
	Version int
}
