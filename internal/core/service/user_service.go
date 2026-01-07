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

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	limit := pageSize

	// Get total count for metadata
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated list
	users, err := s.repo.List(ctx, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
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
