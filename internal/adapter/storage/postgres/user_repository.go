package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gostructure/app/internal/adapter/storage/postgres/dbgen"
	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
	"github.com/gostructure/app/pkg/util/time_util"
)

type UserRepository struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func NewUserRepository(db *sql.DB) port.UserRepository {
	return &UserRepository{
		db:      db,
		queries: dbgen.New(db),
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) List(ctx context.Context, skip, limit int) ([]*domain.User, error) {
	users, err := r.queries.ListUsers(ctx, dbgen.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(skip),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = toDomainUser(u)
	}
	return result, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.TokenVersion == 0 {
		user.TokenVersion = 1
	}
	id, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		Status:       user.Status,
		TokenVersion: int32(user.TokenVersion),
	})
	if err != nil {
		return err
	}

	user.ID = id

	// Fetch the created user to get the database-generated timestamps
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err == nil {
		updatedDomain := toDomainUser(dbUser)
		user.CreatedAt = updatedDomain.CreatedAt
		user.UpdatedAt = updatedDomain.UpdatedAt
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	id, err := r.queries.UpdateUser(ctx, dbgen.UpdateUserParams{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		TokenVersion: int32(user.TokenVersion),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found or no change")
		}
		return err
	}

	if id == 0 {
		return fmt.Errorf("user not found or no change")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// Mapper function: dbgen model to domain model
func toDomainUser(u dbgen.User) *domain.User {
	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		Status:       u.Status,
		TokenVersion: int(u.TokenVersion),
		CreatedAt:    time_util.JSONTime(u.CreatedAt.Time),
		UpdatedAt:    time_util.JSONTime(u.UpdatedAt.Time),
	}
}
