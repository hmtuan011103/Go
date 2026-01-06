package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
)

type InMemoryUserRepository struct {
	users     map[int64]*domain.User
	userMutex sync.RWMutex
	userIDSeq int64
}

func NewUserRepository() port.UserRepository {
	return &InMemoryUserRepository{
		users:     make(map[int64]*domain.User),
		userIDSeq: 1,
	}
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	r.userMutex.RLock()
	defer r.userMutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	r.userMutex.RLock()
	defer r.userMutex.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.userMutex.Lock()
	defer r.userMutex.Unlock()

	user.ID = r.userIDSeq
	r.users[r.userIDSeq] = user
	r.userIDSeq++
	return nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.userMutex.Lock()
	defer r.userMutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id int64) error {
	r.userMutex.Lock()
	defer r.userMutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
