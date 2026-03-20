package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type UserRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.queries.CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           uuidToPgtype(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		DisplayName:  user.DisplayName,
		AvatarUrl:    textToPgtype(user.AvatarURL),
		Bio:          user.Bio,
		CreatedAt:    timestamptzToPgtype(user.CreatedAt),
		UpdatedAt:    timestamptzToPgtype(user.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return sqlcUserToDomain(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return sqlcUserToDomain(row), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return sqlcUserToDomain(row), nil
}

func sqlcUserToDomain(u sqlcgen.User) *domain.User {
	var avatarURL *string
	if u.AvatarUrl.Valid {
		avatarURL = &u.AvatarUrl.String
	}

	return &domain.User{
		ID:           pgtypeToUUID(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
		AvatarURL:    avatarURL,
		Bio:          u.Bio,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func pgtypeToUUID(id pgtype.UUID) uuid.UUID {
	return uuid.UUID(id.Bytes)
}

func textToPgtype(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
