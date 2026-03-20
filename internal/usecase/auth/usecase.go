package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type TokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	Set(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error
	Get(ctx context.Context, userID uuid.UUID) (string, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type RegisterInput struct {
	Username    string
	Email       string
	Password    string
	DisplayName string
}

type LoginInput struct {
	Email    string
	Password string
}

type UseCase struct {
	userRepo     UserRepository
	tokenRepo    TokenRepository
	sessionRepo  SessionRepository
	hasher       PasswordHasher
	tokenManager *jwtpkg.TokenManager
}

func NewUseCase(
	userRepo UserRepository,
	tokenRepo TokenRepository,
	sessionRepo SessionRepository,
	hasher PasswordHasher,
	tokenManager *jwtpkg.TokenManager,
) *UseCase {
	return &UseCase{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		sessionRepo:  sessionRepo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (uc *UseCase) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	_, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, fmt.Errorf("email: %w", domain.ErrAlreadyExists)
	}
	if err != domain.ErrNotFound {
		return nil, fmt.Errorf("check email: %w", err)
	}

	_, err = uc.userRepo.GetByUsername(ctx, input.Username)
	if err == nil {
		return nil, fmt.Errorf("username: %w", domain.ErrAlreadyExists)
	}
	if err != domain.ErrNotFound {
		return nil, fmt.Errorf("check username: %w", err)
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
		DisplayName:  input.DisplayName,
		Bio:          "",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (uc *UseCase) Login(ctx context.Context, input LoginInput) (*jwtpkg.TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if !uc.hasher.Compare(user.PasswordHash, input.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	tokenPair, err := uc.tokenManager.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	refreshHash := hashToken(tokenPair.RefreshToken)
	claims, err := uc.tokenManager.ParseToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	if err := uc.tokenRepo.Create(ctx, user.ID, refreshHash, claims.ExpiresAt.Time); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &tokenPair, nil
}

func (uc *UseCase) Logout(ctx context.Context, refreshToken string) error {
	refreshHash := hashToken(refreshToken)
	if err := uc.tokenRepo.Revoke(ctx, refreshHash); err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

func (uc *UseCase) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if err := uc.tokenRepo.RevokeAllByUser(ctx, userID); err != nil {
		return fmt.Errorf("revoke all tokens: %w", err)
	}
	return nil
}

func (uc *UseCase) RefreshTokens(ctx context.Context, refreshToken string) (*jwtpkg.TokenPair, error) {
	claims, err := uc.tokenManager.ParseToken(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	refreshHash := hashToken(refreshToken)
	storedToken, err := uc.tokenRepo.GetByHash(ctx, refreshHash)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if storedToken.RevokedAt != nil {
		return nil, domain.ErrTokenRevoked
	}

	if err = uc.tokenRepo.Revoke(ctx, refreshHash); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	newPair, err := uc.tokenManager.GenerateTokenPair(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("generate new tokens: %w", err)
	}

	newHash := hashToken(newPair.RefreshToken)
	newClaims, err := uc.tokenManager.ParseToken(newPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse new refresh token: %w", err)
	}

	if err := uc.tokenRepo.Create(ctx, claims.UserID, newHash, newClaims.ExpiresAt.Time); err != nil {
		return nil, fmt.Errorf("save new refresh token: %w", err)
	}

	return &newPair, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
