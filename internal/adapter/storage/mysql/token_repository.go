package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gostructure/app/internal/adapter/storage/mysql/dbgen"
	"github.com/gostructure/app/internal/core/domain"
	"github.com/gostructure/app/internal/core/port"
)

type RefreshTokenRepository struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func NewTokenRepository(db *sql.DB) port.TokenRepository {
	return &RefreshTokenRepository{
		db:      db,
		queries: dbgen.New(db),
	}
}

func (r *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	res, err := r.queries.CreateRefreshToken(ctx, dbgen.CreateRefreshTokenParams{
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}

	id, _ := res.LastInsertId()
	rt.ID = id
	return nil
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	rt, err := r.queries.GetRefreshToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, err
	}

	return toDomainRefreshToken(rt), nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.queries.RevokeRefreshToken(ctx, dbgen.RevokeRefreshTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (r *RefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	return r.queries.RevokeAllUserTokens(ctx, dbgen.RevokeAllUserTokensParams{
		UserID: userID,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func toDomainRefreshToken(rt dbgen.RefreshToken) *domain.RefreshToken {
	var revokedAt *time.Time
	if rt.RevokedAt.Valid {
		revokedAt = &rt.RevokedAt.Time
	}

	return &domain.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt.Time,
		RevokedAt: revokedAt,
	}
}
