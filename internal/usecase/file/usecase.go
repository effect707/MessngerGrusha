package file

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type FileStorage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type AttachmentRepository interface {
	Create(ctx context.Context, a *domain.Attachment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error)
	GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Attachment, error)
}

type MessageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
}

type ChatRepository interface {
	IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

type UseCase struct {
	fileStorage    FileStorage
	attachmentRepo AttachmentRepository
	msgRepo        MessageRepository
	chatRepo       ChatRepository
}

func NewUseCase(
	fileStorage FileStorage,
	attachmentRepo AttachmentRepository,
	msgRepo MessageRepository,
	chatRepo ChatRepository,
) *UseCase {
	return &UseCase{
		fileStorage:    fileStorage,
		attachmentRepo: attachmentRepo,
		msgRepo:        msgRepo,
		chatRepo:       chatRepo,
	}
}

type UploadInput struct {
	MessageID  uuid.UUID
	UserID     uuid.UUID
	FileName   string
	FileSize   int64
	MimeType   string
	DurationMs *int
	Reader     io.Reader
}

func (uc *UseCase) Upload(ctx context.Context, input UploadInput) (*domain.Attachment, error) {
	msg, err := uc.msgRepo.GetByID(ctx, input.MessageID)
	if err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}

	if msg.SenderID != input.UserID {
		return nil, fmt.Errorf("only sender can attach files: %w", domain.ErrForbidden)
	}

	storageKey := fmt.Sprintf("%s/%s/%s", msg.ChatID, input.MessageID, input.FileName)

	if err := uc.fileStorage.Upload(ctx, storageKey, input.Reader, input.FileSize, input.MimeType); err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	attachment := &domain.Attachment{
		ID:         uuid.New(),
		MessageID:  input.MessageID,
		FileName:   input.FileName,
		FileSize:   input.FileSize,
		MimeType:   input.MimeType,
		StorageKey: storageKey,
		DurationMs: input.DurationMs,
		CreatedAt:  time.Now(),
	}

	if err := uc.attachmentRepo.Create(ctx, attachment); err != nil {
		_ = uc.fileStorage.Delete(ctx, storageKey)
		return nil, fmt.Errorf("save attachment: %w", err)
	}

	return attachment, nil
}

func (uc *UseCase) Download(ctx context.Context, attachmentID, userID uuid.UUID) (io.ReadCloser, *domain.Attachment, error) {
	attachment, err := uc.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("get attachment: %w", err)
	}

	msg, err := uc.msgRepo.GetByID(ctx, attachment.MessageID)
	if err != nil {
		return nil, nil, fmt.Errorf("get message: %w", err)
	}

	if msg.ChatID != nil {
		var isMember bool
		isMember, err = uc.chatRepo.IsMember(ctx, *msg.ChatID, userID)
		if err != nil {
			return nil, nil, fmt.Errorf("check membership: %w", err)
		}
		if !isMember {
			return nil, nil, domain.ErrForbidden
		}
	}

	reader, err := uc.fileStorage.Download(ctx, attachment.StorageKey)
	if err != nil {
		return nil, nil, fmt.Errorf("download file: %w", err)
	}

	return reader, attachment, nil
}

func (uc *UseCase) GetAttachments(ctx context.Context, messageID, userID uuid.UUID) ([]domain.Attachment, error) {
	msg, err := uc.msgRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}

	if msg.ChatID != nil {
		var isMember bool
		isMember, err = uc.chatRepo.IsMember(ctx, *msg.ChatID, userID)
		if err != nil {
			return nil, fmt.Errorf("check membership: %w", err)
		}
		if !isMember {
			return nil, domain.ErrForbidden
		}
	}

	attachments, err := uc.attachmentRepo.GetByMessageID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get attachments: %w", err)
	}
	return attachments, nil
}
