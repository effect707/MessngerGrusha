package notification

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type mockNotifRepo struct {
	mock.Mock
}

func (m *mockNotifRepo) Create(ctx context.Context, n *domain.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *mockNotifRepo) GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *mockNotifRepo) GetUnread(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *mockNotifRepo) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockNotifRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockNotifRepo) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func newTestUseCase() (*UseCase, *mockNotifRepo) {
	repo := new(mockNotifRepo)
	uc := NewUseCase(repo)
	return uc, repo
}

func TestSend_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	err := uc.Send(ctx, userID, domain.NotificationTypeNewMessage, map[string]string{"chat_id": "123"})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetNotifications_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	now := time.Now()
	notifications := []domain.Notification{
		{ID: uuid.New(), UserID: userID, Type: domain.NotificationTypeNewMessage, CreatedAt: now},
		{ID: uuid.New(), UserID: userID, Type: domain.NotificationTypeReaction, CreatedAt: now},
	}
	repo.On("GetByUser", ctx, userID, 50).Return(notifications, nil)

	result, err := uc.GetNotifications(ctx, userID, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetUnread_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	repo.On("GetUnread", ctx, userID, 20).Return([]domain.Notification{
		{ID: uuid.New(), IsRead: false},
	}, nil)

	result, err := uc.GetUnread(ctx, userID, 20)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestMarkRead_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	notifID := uuid.New()
	userID := uuid.New()

	repo.On("MarkRead", ctx, notifID, userID).Return(nil)

	err := uc.MarkRead(ctx, notifID, userID)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkAllRead_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	repo.On("MarkAllRead", ctx, userID).Return(nil)

	err := uc.MarkAllRead(ctx, userID)

	require.NoError(t, err)
}

func TestCountUnread_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	repo.On("CountUnread", ctx, userID).Return(int64(5), nil)

	count, err := uc.CountUnread(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
