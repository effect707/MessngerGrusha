package file

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type mockFileStorage struct {
	mock.Mock
}

func (m *mockFileStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *mockFileStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockFileStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type mockAttachmentRepo struct {
	mock.Mock
}

func (m *mockAttachmentRepo) Create(ctx context.Context, a *domain.Attachment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *mockAttachmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attachment), args.Error(1)
}

func (m *mockAttachmentRepo) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Attachment, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Attachment), args.Error(1)
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

func newTestUseCase() (*UseCase, *mockFileStorage, *mockAttachmentRepo, *mockMsgRepo, *mockChatRepo) {
	fs := new(mockFileStorage)
	ar := new(mockAttachmentRepo)
	mr := new(mockMsgRepo)
	cr := new(mockChatRepo)
	uc := NewUseCase(fs, ar, mr, cr)
	return uc, fs, ar, mr, cr
}

func TestUpload_Success(t *testing.T) {
	uc, fileStorage, attachmentRepo, msgRepo, _ := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()
	content := []byte("file content")
	reader := bytes.NewReader(content)

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:       messageID,
		ChatID:   &chatID,
		SenderID: userID,
	}, nil)
	fileStorage.On("Upload", ctx, mock.AnythingOfType("string"), reader, int64(len(content)), "image/png").Return(nil)
	attachmentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Attachment")).Return(nil)

	attachment, err := uc.Upload(ctx, UploadInput{
		MessageID: messageID,
		UserID:    userID,
		FileName:  "test.png",
		FileSize:  int64(len(content)),
		MimeType:  "image/png",
		Reader:    reader,
	})

	require.NoError(t, err)
	assert.Equal(t, "test.png", attachment.FileName)
	assert.Equal(t, int64(len(content)), attachment.FileSize)
	assert.Equal(t, "image/png", attachment.MimeType)
}

func TestUpload_NotSender(t *testing.T) {
	uc, _, _, msgRepo, _ := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	chatID := uuid.New()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:       messageID,
		ChatID:   &chatID,
		SenderID: otherUserID,
	}, nil)

	_, err := uc.Upload(ctx, UploadInput{
		MessageID: messageID,
		UserID:    userID,
		FileName:  "test.png",
		FileSize:  100,
		MimeType:  "image/png",
		Reader:    bytes.NewReader([]byte("data")),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestDownload_Success(t *testing.T) {
	uc, fileStorage, attachmentRepo, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	attachmentID := uuid.New()
	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	attachmentRepo.On("GetByID", ctx, attachmentID).Return(&domain.Attachment{
		ID:         attachmentID,
		MessageID:  messageID,
		FileName:   "test.png",
		FileSize:   100,
		MimeType:   "image/png",
		StorageKey: "key/test.png",
	}, nil)
	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	fileStorage.On("Download", ctx, "key/test.png").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)

	reader, attachment, err := uc.Download(ctx, attachmentID, userID)

	require.NoError(t, err)
	assert.Equal(t, "test.png", attachment.FileName)
	defer func() { _ = reader.Close() }()
}

func TestDownload_NotMember(t *testing.T) {
	uc, _, attachmentRepo, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	attachmentID := uuid.New()
	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	attachmentRepo.On("GetByID", ctx, attachmentID).Return(&domain.Attachment{
		ID:        attachmentID,
		MessageID: messageID,
	}, nil)
	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(false, nil)

	_, _, err := uc.Download(ctx, attachmentID, userID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestGetAttachments_Success(t *testing.T) {
	uc, _, attachmentRepo, msgRepo, chatRepo := newTestUseCase()
	ctx := context.Background()

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()

	msgRepo.On("GetByID", ctx, messageID).Return(&domain.Message{
		ID:     messageID,
		ChatID: &chatID,
	}, nil)
	chatRepo.On("IsMember", ctx, chatID, userID).Return(true, nil)
	attachmentRepo.On("GetByMessageID", ctx, messageID).Return([]domain.Attachment{
		{ID: uuid.New(), FileName: "a.png"},
		{ID: uuid.New(), FileName: "b.pdf"},
	}, nil)

	attachments, err := uc.GetAttachments(ctx, messageID, userID)

	require.NoError(t, err)
	assert.Len(t, attachments, 2)
}
