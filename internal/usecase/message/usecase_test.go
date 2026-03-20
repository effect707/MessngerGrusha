package message

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/pkg/pagination"
)

type mockMessageRepo struct {
	mock.Mock
}

func (m *mockMessageRepo) Create(ctx context.Context, msg *domain.Message) error {
	return m.Called(ctx, msg).Error(0)
}

func (m *mockMessageRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *mockMessageRepo) GetChatHistory(ctx context.Context, chatID uuid.UUID, cursor *pagination.Cursor, limit int) ([]domain.Message, error) {
	args := m.Called(ctx, chatID, cursor, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *mockMessageRepo) Search(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]domain.Message, error) {
	args := m.Called(ctx, chatID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *mockMessageRepo) Update(ctx context.Context, msg *domain.Message) error {
	return m.Called(ctx, msg).Error(0)
}

func (m *mockMessageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockChatRepo struct {
	mock.Mock
}

func (m *mockChatRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *mockChatRepo) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, chatID, userID)
	return args.Bool(0), args.Error(1)
}

type mockBroker struct {
	mock.Mock
}

func (m *mockBroker) Publish(ctx context.Context, chatID string, msg []byte) error {
	return m.Called(ctx, chatID, msg).Error(0)
}

func newTestUseCase() (*UseCase, *mockMessageRepo, *mockChatRepo, *mockBroker) {
	msgRepo := new(mockMessageRepo)
	chatRepo := new(mockChatRepo)
	broker := new(mockBroker)

	uc := NewUseCase(msgRepo, chatRepo, broker)
	return uc, msgRepo, chatRepo, broker
}

func TestSend_Success(t *testing.T) {
	uc, msgRepo, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	msgRepo.On("Create", ctx, mock.AnythingOfType("*domain.Message")).Return(nil)

	msg, err := uc.Send(ctx, SendInput{
		ChatID:   chatID,
		SenderID: userID,
		Type:     domain.MessageTypeText,
		Content:  "Hello, World!",
	})

	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", msg.Content)
	assert.Equal(t, domain.MessageTypeText, msg.Type)
	assert.Equal(t, userID, msg.SenderID)
	assert.Equal(t, &chatID, msg.ChatID)
}

func TestSend_NotMember(t *testing.T) {
	uc, _, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	_, err := uc.Send(ctx, SendInput{
		ChatID:   chatID,
		SenderID: userID,
		Type:     domain.MessageTypeText,
		Content:  "Hello!",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestSend_EmptyTextMessage(t *testing.T) {
	uc, _, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)

	_, err := uc.Send(ctx, SendInput{
		ChatID:   chatID,
		SenderID: userID,
		Type:     domain.MessageTypeText,
		Content:  "",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestGetHistory_Success(t *testing.T) {
	uc, msgRepo, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	messages := []domain.Message{
		{ID: uuid.New(), Content: "msg1", CreatedAt: time.Now()},
		{ID: uuid.New(), Content: "msg2", CreatedAt: time.Now().Add(-time.Minute)},
	}

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	msgRepo.On("GetChatHistory", ctx, chatID, (*pagination.Cursor)(nil), 51).Return(messages, nil)

	page, err := uc.GetHistory(ctx, chatID, userID, nil, 50)

	require.NoError(t, err)
	assert.Len(t, page.Items, 2)
	assert.False(t, page.HasMore)
	assert.Nil(t, page.NextCursor)
}

func TestGetHistory_WithPagination(t *testing.T) {
	uc, msgRepo, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	now := time.Now()
	messages := make([]domain.Message, 51)
	for i := range messages {
		messages[i] = domain.Message{
			ID:        uuid.New(),
			Content:   "msg",
			CreatedAt: now.Add(-time.Duration(i) * time.Minute),
		}
	}

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	msgRepo.On("GetChatHistory", ctx, chatID, (*pagination.Cursor)(nil), 51).Return(messages, nil)

	page, err := uc.GetHistory(ctx, chatID, userID, nil, 50)

	require.NoError(t, err)
	assert.Len(t, page.Items, 50)
	assert.True(t, page.HasMore)
	assert.NotNil(t, page.NextCursor)
}

func TestSearchMessages_Success(t *testing.T) {
	uc, msgRepo, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	messages := []domain.Message{
		{ID: uuid.New(), Content: "hello world"},
	}

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	msgRepo.On("Search", ctx, chatID, "hello", 50).Return(messages, nil)

	result, err := uc.SearchMessages(ctx, chatID, userID, "hello", 50)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}
