package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/effect707/MessngerGrusha/internal/domain"
	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type mockTokenRepo struct {
	mock.Mock
}

func (m *mockTokenRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *mockTokenRepo) GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshToken), args.Error(1)
}

func (m *mockTokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *mockTokenRepo) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) Set(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	args := m.Called(ctx, userID, token, ttl)
	return args.Error(0)
}

func (m *mockSessionRepo) Get(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockSessionRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockHasher) Compare(hash, password string) bool {
	return hash == "hashed_"+password
}

func newTestUseCase() (*UseCase, *mockUserRepo, *mockTokenRepo, *mockSessionRepo) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	sessionRepo := new(mockSessionRepo)
	hasher := &mockHasher{}
	tm := jwtpkg.NewTokenManager("test-secret", 15*time.Minute, 720*time.Hour)

	uc := NewUseCase(userRepo, tokenRepo, sessionRepo, hasher, tm)
	return uc, userRepo, tokenRepo, sessionRepo
}

func TestRegister_Success(t *testing.T) {
	uc, userRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, domain.ErrNotFound)
	userRepo.On("GetByUsername", ctx, "testuser").Return(nil, domain.ErrNotFound)
	userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := uc.Register(ctx, RegisterInput{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "hashed_password123", user.PasswordHash)
	assert.NotEqual(t, uuid.Nil, user.ID)

	userRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	uc, userRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	existingUser := &domain.User{Email: "test@example.com"}
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

	_, err := uc.Register(ctx, RegisterInput{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	uc, userRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, domain.ErrNotFound)
	existingUser := &domain.User{Username: "testuser"}
	userRepo.On("GetByUsername", ctx, "testuser").Return(existingUser, nil)

	_, err := uc.Register(ctx, RegisterInput{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestLogin_Success(t *testing.T) {
	uc, userRepo, tokenRepo, _ := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
	}
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	tokenRepo.On("Create", ctx, userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	tokens, err := uc.Login(ctx, LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	uc, userRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
	}
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

	_, err := uc.Login(ctx, LoginInput{
		Email:    "test@example.com",
		Password: "wrong_password",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, userRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, domain.ErrNotFound)

	_, err := uc.Login(ctx, LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogout_Success(t *testing.T) {
	uc, _, tokenRepo, _ := newTestUseCase()
	ctx := context.Background()

	tokenRepo.On("Revoke", ctx, mock.AnythingOfType("string")).Return(nil)

	err := uc.Logout(ctx, "some-refresh-token")
	require.NoError(t, err)

	tokenRepo.AssertExpectations(t)
}

func TestLogoutAll_Success(t *testing.T) {
	uc, _, tokenRepo, _ := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	tokenRepo.On("RevokeAllByUser", ctx, userID).Return(nil)

	err := uc.LogoutAll(ctx, userID)
	require.NoError(t, err)

	tokenRepo.AssertExpectations(t)
}
