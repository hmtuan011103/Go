package service

import (
	"context"
	"errors"

	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	// Business logic could be here: validate email format, check uniqueness, etc.

	user := &domain.User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
