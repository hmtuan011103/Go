package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
)

type AuthService struct {
	userRepo  port.UserRepository
	tokenRepo port.TokenRepository
	hasher    port.PasswordHasher
	tokenProv port.TokenProvider
}

func NewAuthService(
	userRepo port.UserRepository,
	tokenRepo port.TokenRepository,
	hasher port.PasswordHasher,
	tokenProv port.TokenProvider,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		hasher:    hasher,
		tokenProv: tokenProv,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email and password are required")
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*port.AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := s.hasher.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 1. INCREMENT TOKEN VERSION TO INVALIDATE ALL OLD JWTS
	user.TokenVersion++
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update token version: %w", err)
	}

	// 2. INVALIDATE ALL OLD REFRESH TOKENS FOR THIS USER
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old tokens: %w", err)
	}

	// 3. Create Access Token (JWT - short-lived)
	accessToken, err := s.tokenProv.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 4. Create Refresh Token (Random string - long-lived)
	refreshTokenStr, err := s.tokenProv.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Save Refresh Token to Database
	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	if err := s.tokenRepo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &port.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshTokenStr string) (*port.AuthResult, error) {
	// 1. Get token from DB
	rt, err := s.tokenRepo.GetRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Validate token
	if rt.IsRevoked() || rt.IsExpired() {
		return nil, errors.New("refresh token expired or revoked")
	}

	// 3. Get user
	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. Create new Access Token
	accessToken, err := s.tokenProv.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &port.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	return s.tokenRepo.RevokeRefreshToken(ctx, refreshTokenStr)
}
