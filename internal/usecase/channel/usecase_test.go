package channel

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

type mockChannelRepo struct {
	mock.Mock
}

func (m *mockChannelRepo) Create(ctx context.Context, ch *domain.Channel) error {
	args := m.Called(ctx, ch)
	return args.Error(0)
}

func (m *mockChannelRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Channel), args.Error(1)
}

func (m *mockChannelRepo) GetBySlug(ctx context.Context, slug string) (*domain.Channel, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Channel), args.Error(1)
}

func (m *mockChannelRepo) Update(ctx context.Context, ch *domain.Channel) error {
	args := m.Called(ctx, ch)
	return args.Error(0)
}

func (m *mockChannelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockChannelRepo) GetPublic(ctx context.Context, limit int) ([]domain.Channel, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Channel), args.Error(1)
}

func (m *mockChannelRepo) GetUserChannels(ctx context.Context, userID uuid.UUID) ([]domain.Channel, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Channel), args.Error(1)
}

func (m *mockChannelRepo) AddSubscriber(ctx context.Context, sub *domain.ChannelSubscriber) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *mockChannelRepo) RemoveSubscriber(ctx context.Context, channelID, userID uuid.UUID) error {
	args := m.Called(ctx, channelID, userID)
	return args.Error(0)
}

func (m *mockChannelRepo) IsSubscriber(ctx context.Context, channelID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, channelID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockChannelRepo) GetSubscribers(ctx context.Context, channelID uuid.UUID) ([]domain.ChannelSubscriber, error) {
	args := m.Called(ctx, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ChannelSubscriber), args.Error(1)
}

func newTestUseCase() (*UseCase, *mockChannelRepo) {
	repo := new(mockChannelRepo)
	uc := NewUseCase(repo)
	return uc, repo
}

func TestCreate_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()
	ownerID := uuid.New()

	repo.On("GetBySlug", ctx, "test-channel").Return(nil, domain.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Channel")).Return(nil)
	repo.On("AddSubscriber", ctx, mock.AnythingOfType("*domain.ChannelSubscriber")).Return(nil)

	ch, err := uc.Create(ctx, CreateInput{
		Slug:    "test-channel",
		Name:    "Test Channel",
		OwnerID: ownerID,
	})

	require.NoError(t, err)
	assert.Equal(t, "test-channel", ch.Slug)
	assert.Equal(t, "Test Channel", ch.Name)
	assert.Equal(t, ownerID, ch.OwnerID)
	repo.AssertExpectations(t)
}

func TestCreate_EmptySlug(t *testing.T) {
	uc, _ := newTestUseCase()
	ctx := context.Background()

	_, err := uc.Create(ctx, CreateInput{Name: "Test"})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreate_EmptyName(t *testing.T) {
	uc, _ := newTestUseCase()
	ctx := context.Background()

	_, err := uc.Create(ctx, CreateInput{Slug: "test"})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreate_SlugTaken(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	repo.On("GetBySlug", ctx, "taken").Return(&domain.Channel{Slug: "taken"}, nil)

	_, err := uc.Create(ctx, CreateInput{Slug: "taken", Name: "Test"})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestGetByID_PrivateChannel_NotSubscriber(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	userID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:        channelID,
		IsPrivate: true,
	}, nil)
	repo.On("IsSubscriber", ctx, channelID, userID).Return(false, nil)

	_, err := uc.GetByID(ctx, channelID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestGetByID_PublicChannel(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	userID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:        channelID,
		IsPrivate: false,
	}, nil)

	ch, err := uc.GetByID(ctx, channelID, userID)

	require.NoError(t, err)
	assert.Equal(t, channelID, ch.ID)
}

func TestUpdate_NotOwner(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	ownerID := uuid.New()
	otherID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:      channelID,
		OwnerID: ownerID,
	}, nil)

	_, err := uc.Update(ctx, channelID, otherID, UpdateInput{Name: "New Name"})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestUpdate_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	ownerID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:      channelID,
		OwnerID: ownerID,
		Name:    "Old",
	}, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*domain.Channel")).Return(nil)

	ch, err := uc.Update(ctx, channelID, ownerID, UpdateInput{Name: "New Name", Description: "desc"})

	require.NoError(t, err)
	assert.Equal(t, "New Name", ch.Name)
}

func TestDelete_NotOwner(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	ownerID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:      channelID,
		OwnerID: ownerID,
	}, nil)

	err := uc.Delete(ctx, channelID, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestSubscribe_PrivateChannel(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:        channelID,
		IsPrivate: true,
	}, nil)

	err := uc.Subscribe(ctx, channelID, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestSubscribe_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	userID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:        channelID,
		IsPrivate: false,
	}, nil)
	repo.On("AddSubscriber", ctx, mock.AnythingOfType("*domain.ChannelSubscriber")).Return(nil)

	err := uc.Subscribe(ctx, channelID, userID)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUnsubscribe_OwnerCannot(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	channelID := uuid.New()
	ownerID := uuid.New()

	repo.On("GetByID", ctx, channelID).Return(&domain.Channel{
		ID:      channelID,
		OwnerID: ownerID,
	}, nil)

	err := uc.Unsubscribe(ctx, channelID, ownerID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestGetPublicChannels_Success(t *testing.T) {
	uc, repo := newTestUseCase()
	ctx := context.Background()

	now := time.Now()
	channels := []domain.Channel{
		{ID: uuid.New(), Name: "Ch1", CreatedAt: now},
		{ID: uuid.New(), Name: "Ch2", CreatedAt: now},
	}
	repo.On("GetPublic", ctx, 50).Return(channels, nil)

	result, err := uc.GetPublicChannels(ctx, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
