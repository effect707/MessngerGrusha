package chat

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

func (m *mockChatRepo) CreateDirect(ctx context.Context, creatorID, recipientID uuid.UUID) (*domain.Chat, error) {
	args := m.Called(ctx, creatorID, recipientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *mockChatRepo) CreateGroup(ctx context.Context, name string, creatorID uuid.UUID, memberIDs []uuid.UUID) (*domain.Chat, error) {
	args := m.Called(ctx, name, creatorID, memberIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *mockChatRepo) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *mockChatRepo) GetDirectChatID(ctx context.Context, userA, userB uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, userA, userB)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockChatRepo) AddMember(ctx context.Context, chatID, userID uuid.UUID, role domain.MemberRole) error {
	args := m.Called(ctx, chatID, userID, role)
	return args.Error(0)
}

func (m *mockChatRepo) RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error {
	args := m.Called(ctx, chatID, userID)
	return args.Error(0)
}

func (m *mockChatRepo) GetMembers(ctx context.Context, chatID uuid.UUID) ([]domain.ChatMember, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ChatMember), args.Error(1)
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func newTestUseCase() (*UseCase, *mockChatRepo, *mockUserRepo) {
	chatRepo := new(mockChatRepo)
	userRepo := new(mockUserRepo)
	uc := NewUseCase(chatRepo, userRepo)
	return uc, chatRepo, userRepo
}

func TestCreateDirectChat_Success(t *testing.T) {
	uc, chatRepo, userRepo := newTestUseCase()
	ctx := context.Background()

	creatorID := uuid.New()
	recipientID := uuid.New()
	chatID := uuid.New()
	now := time.Now()

	userRepo.On("GetByID", ctx, recipientID).Return(&domain.User{ID: recipientID}, nil)
	chatRepo.On("GetDirectChatID", ctx, creatorID, recipientID).Return(uuid.UUID{}, domain.ErrNotFound)
	chatRepo.On("CreateDirect", ctx, creatorID, recipientID).Return(&domain.Chat{
		ID:        chatID,
		Type:      domain.ChatTypeDirect,
		CreatedBy: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	chat, err := uc.CreateDirectChat(ctx, creatorID, recipientID)

	require.NoError(t, err)
	assert.Equal(t, chatID, chat.ID)
	assert.Equal(t, domain.ChatTypeDirect, chat.Type)
	chatRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCreateDirectChat_SelfChat(t *testing.T) {
	uc, _, _ := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()

	_, err := uc.CreateDirectChat(ctx, userID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreateDirectChat_ExistingChat(t *testing.T) {
	uc, chatRepo, userRepo := newTestUseCase()
	ctx := context.Background()

	creatorID := uuid.New()
	recipientID := uuid.New()
	existingChatID := uuid.New()
	now := time.Now()

	existingChat := &domain.Chat{
		ID:        existingChatID,
		Type:      domain.ChatTypeDirect,
		CreatedBy: creatorID,
		CreatedAt: now,
	}

	userRepo.On("GetByID", ctx, recipientID).Return(&domain.User{ID: recipientID}, nil)
	chatRepo.On("GetDirectChatID", ctx, creatorID, recipientID).Return(existingChatID, nil)
	chatRepo.On("GetByID", ctx, existingChatID).Return(existingChat, nil)

	chat, err := uc.CreateDirectChat(ctx, creatorID, recipientID)

	require.NoError(t, err)
	assert.Equal(t, existingChatID, chat.ID)
	chatRepo.AssertNotCalled(t, "CreateDirect")
}

func TestCreateDirectChat_RecipientNotFound(t *testing.T) {
	uc, _, userRepo := newTestUseCase()
	ctx := context.Background()

	creatorID := uuid.New()
	recipientID := uuid.New()

	userRepo.On("GetByID", ctx, recipientID).Return(nil, domain.ErrNotFound)

	_, err := uc.CreateDirectChat(ctx, creatorID, recipientID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCreateGroupChat_Success(t *testing.T) {
	uc, chatRepo, userRepo := newTestUseCase()
	ctx := context.Background()

	creatorID := uuid.New()
	member1 := uuid.New()
	member2 := uuid.New()
	chatID := uuid.New()
	name := "Test Group"
	now := time.Now()

	userRepo.On("GetByID", ctx, member1).Return(&domain.User{ID: member1}, nil)
	userRepo.On("GetByID", ctx, member2).Return(&domain.User{ID: member2}, nil)
	chatRepo.On("CreateGroup", ctx, name, creatorID, []uuid.UUID{member1, member2}).Return(&domain.Chat{
		ID:        chatID,
		Type:      domain.ChatTypeGroup,
		Name:      &name,
		CreatedBy: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	chat, err := uc.CreateGroupChat(ctx, name, creatorID, []uuid.UUID{member1, member2})

	require.NoError(t, err)
	assert.Equal(t, chatID, chat.ID)
	assert.Equal(t, domain.ChatTypeGroup, chat.Type)
	assert.Equal(t, &name, chat.Name)
	chatRepo.AssertExpectations(t)
}

func TestCreateGroupChat_EmptyName(t *testing.T) {
	uc, _, _ := newTestUseCase()
	ctx := context.Background()

	_, err := uc.CreateGroupChat(ctx, "", uuid.New(), []uuid.UUID{uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreateGroupChat_MemberNotFound(t *testing.T) {
	uc, _, userRepo := newTestUseCase()
	ctx := context.Background()

	memberID := uuid.New()
	userRepo.On("GetByID", ctx, memberID).Return(nil, domain.ErrNotFound)

	_, err := uc.CreateGroupChat(ctx, "Test", uuid.New(), []uuid.UUID{memberID})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGetChat_Success(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()
	chat := &domain.Chat{ID: chatID, Type: domain.ChatTypeDirect}

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)

	result, err := uc.GetChat(ctx, chatID, userID)

	require.NoError(t, err)
	assert.Equal(t, chatID, result.ID)
}

func TestGetChat_NotMember(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	_, err := uc.GetChat(ctx, chatID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestGetUserChats_Success(t *testing.T) {
	uc, chatRepo, userRepo := newTestUseCase()
	ctx := context.Background()

	userID := uuid.New()
	otherUserID := uuid.New()
	directChatID := uuid.New()
	groupChatID := uuid.New()

	chats := []domain.Chat{
		{ID: directChatID, Type: domain.ChatTypeDirect},
		{ID: groupChatID, Type: domain.ChatTypeGroup},
	}

	chatRepo.On("GetUserChats", ctx, userID).Return(chats, nil)
	chatRepo.On("GetMembers", ctx, directChatID).Return([]domain.ChatMember{
		{ChatID: directChatID, UserID: userID, Role: domain.MemberRoleMember},
		{ChatID: directChatID, UserID: otherUserID, Role: domain.MemberRoleMember},
	}, nil)
	userRepo.On("GetByID", ctx, otherUserID).Return(&domain.User{
		ID:       otherUserID,
		Username: "other_user",
	}, nil)

	result, err := uc.GetUserChats(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "other_user", *result[0].Name)
}

func TestAddMember_Success(t *testing.T) {
	uc, chatRepo, userRepo := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	adderID := uuid.New()
	newMemberID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeGroup,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, adderID).Return(true, nil)
	chatRepo.On("IsMember", ctx, chatID, newMemberID).Return(false, nil)
	userRepo.On("GetByID", ctx, newMemberID).Return(&domain.User{ID: newMemberID}, nil)
	chatRepo.On("AddMember", ctx, chatID, newMemberID, domain.MemberRoleMember).Return(nil)

	err := uc.AddMember(ctx, chatID, adderID, newMemberID)

	require.NoError(t, err)
	chatRepo.AssertExpectations(t)
}

func TestAddMember_DirectChat(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeDirect,
	}, nil)

	err := uc.AddMember(ctx, chatID, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAddMember_NotMember(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	adderID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeGroup,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, adderID).Return(false, nil)

	err := uc.AddMember(ctx, chatID, adderID, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestAddMember_AlreadyMember(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	adderID := uuid.New()
	newMemberID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeGroup,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, adderID).Return(true, nil)
	chatRepo.On("IsMember", ctx, chatID, newMemberID).Return(true, nil)

	err := uc.AddMember(ctx, chatID, adderID, newMemberID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestRemoveMember_Success(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	removerID := uuid.New()
	targetID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeGroup,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, removerID).Return(true, nil)
	chatRepo.On("RemoveMember", ctx, chatID, targetID).Return(nil)

	err := uc.RemoveMember(ctx, chatID, removerID, targetID)

	require.NoError(t, err)
	chatRepo.AssertExpectations(t)
}

func TestRemoveMember_NotMember(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	removerID := uuid.New()

	chatRepo.On("GetByID", ctx, chatID).Return(&domain.Chat{
		ID:   chatID,
		Type: domain.ChatTypeGroup,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, removerID).Return(false, nil)

	err := uc.RemoveMember(ctx, chatID, removerID, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestGetMembers_Success(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	members := []domain.ChatMember{
		{ChatID: chatID, UserID: userID, Role: domain.MemberRoleAdmin, JoinedAt: now},
		{ChatID: chatID, UserID: uuid.New(), Role: domain.MemberRoleMember, JoinedAt: now},
	}

	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	chatRepo.On("GetMembers", ctx, chatID).Return(members, nil)

	result, err := uc.GetMembers(ctx, chatID, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetMembers_NotMember(t *testing.T) {
	uc, chatRepo, _ := newTestUseCase()
	ctx := context.Background()

	chatID := uuid.New()
	userID := uuid.New()

	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	_, err := uc.GetMembers(ctx, chatID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}
