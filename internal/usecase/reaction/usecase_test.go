package reaction

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

type mockReactionRepo struct {
	mock.Mock
}

func (m *mockReactionRepo) Add(ctx context.Context, reaction *domain.Reaction) error {
	args := m.Called(ctx, reaction)
	return args.Error(0)
}

func (m *mockReactionRepo) Remove(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	args := m.Called(ctx, messageID, userID, emoji)
	return args.Error(0)
}

func (m *mockReactionRepo) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Reaction, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Reaction), args.Error(1)
}

type mockMsgRepo struct {
	mock.Mock
}

func (m *mockMsgRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

type mockChatRepo struct {
	mock.Mock
}

func (m *mockChatRepo) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, chatID, userID)
	return args.Bool(0), args.Error(1)
}

func newTestUseCase() (*UseCase, *mockReactionRepo, *mockMsgRepo, *mockChatRepo) {
	reactionRepo := new(mockReactionRepo)
	msgRepo := new(mockMsgRepo)
	chatRepo := new(mockChatRepo)
	uc := NewUseCase(reactionRepo, msgRepo, chatRepo)
	return uc, reactionRepo, msgRepo, chatRepo
}

func TestAddReaction_Success(t *testing.T) {
	uc, reactionRepo, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	reactionRepo.On("Add", ctx, mock.AnythingOfType("*domain.Reaction")).Return(nil)

	err := uc.AddReaction(ctx, messageID, userID, "👍")

	require.NoError(t, err)
	reactionRepo.AssertExpectations(t)
}

func TestAddReaction_EmptyEmoji(t *testing.T) {
	uc, _, _, _ := newTestUseCase()
	ctx := context.Background()

	err := uc.AddReaction(ctx, uuid.New(), uuid.New(), "")

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAddReaction_NotMember(t *testing.T) {
	uc, _, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	err := uc.AddReaction(ctx, messageID, userID, "👍")

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestAddReaction_MessageNotFound(t *testing.T) {
	uc, _, msgRepo, _ := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	msgRepo.On("GetByID", ctx, messageID).Return(nil, domain.ErrNotFound)

	err := uc.AddReaction(ctx, messageID, uuid.New(), "👍")

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRemoveReaction_Success(t *testing.T) {
	uc, reactionRepo, _, _ := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()

	reactionRepo.On("Remove", ctx, messageID, userID, "👍").Return(nil)

	err := uc.RemoveReaction(ctx, messageID, userID, "👍")

	require.NoError(t, err)
	reactionRepo.AssertExpectations(t)
}

func TestGetReactions_Success(t *testing.T) {
	uc, reactionRepo, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()
	now := time.Now()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	reactionRepo.On("GetByMessageID", ctx, messageID).Return([]domain.Reaction{
		{MessageID: messageID, UserID: uuid.New(), Emoji: "👍", CreatedAt: now},
		{MessageID: messageID, UserID: uuid.New(), Emoji: "❤️", CreatedAt: now},
	}, nil)

	reactions, err := uc.GetReactions(ctx, messageID, userID)

	require.NoError(t, err)
	assert.Len(t, reactions, 2)
}

func TestGetReactions_NotMember(t *testing.T) {
	uc, _, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	_, err := uc.GetReactions(ctx, messageID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}
