package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gostructure/app/internal/config"
	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
)

type JWTProvider struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTProvider(cfg *config.JWTConfig) *JWTProvider {
	return &JWTProvider{
		secret:     []byte(cfg.Secret),
		expiration: cfg.Expiration, // Usually 15m in .env
	}
}

func (p *JWTProvider) GenerateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"version": user.TokenVersion,
		"exp":     time.Now().Add(p.expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JWTProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (p *JWTProvider) ValidateToken(tokenString string) (*port.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 1. Extract UserID
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid or missing user_id in token")
		}

		// 2. Extract Role
		role, ok := claims["role"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid or missing role in token")
		}

		// 3. Extract Version
		versionFloat, ok := claims["version"].(float64)
		if !ok {
			return nil, fmt.Errorf("token is too old and lacks security version")
		}

		return &port.Claims{
			UserID:  int64(userIDFloat),
			Role:    role,
			Version: int(versionFloat),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
