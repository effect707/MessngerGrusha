package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
	authuc "github.com/effect707/MessngerGrusha/internal/usecase/auth"
)

type TokenRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *TokenRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	err := r.queries.CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		UserID:    uuidToPgtype(userID),
		TokenHash: tokenHash,
		ExpiresAt: timestamptzToPgtype(expiresAt),
	})
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepository) GetByHash(ctx context.Context, tokenHash string) (*authuc.RefreshToken, error) {
	row, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	token := &authuc.RefreshToken{
		ID:        pgtypeToUUID(row.ID),
		UserID:    pgtypeToUUID(row.UserID),
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.RevokedAt.Valid {
		t := row.RevokedAt.Time
		token.RevokedAt = &t
	}

	return token, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	err := r.queries.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.RevokeAllUserRefreshTokens(ctx, uuidToPgtype(userID))
	if err != nil {
		return fmt.Errorf("revoke all tokens: %w", err)
	}
	return nil
}
