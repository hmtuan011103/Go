package port

import (
	"context"

	"github.com/gostructure/app/internal/core/domain"
)

// UserRepository defines the interface for data access
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, skip, limit int) ([]*domain.User, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

// UserService defines the interface for business logic
type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error)
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
