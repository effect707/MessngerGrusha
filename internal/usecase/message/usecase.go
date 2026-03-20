package message

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/pkg/pagination"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	GetChatHistory(ctx context.Context, chatID uuid.UUID, cursor *pagination.Cursor, limit int) ([]domain.Message, error)
	Search(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]domain.Message, error)
	Update(ctx context.Context, msg *domain.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChatRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

type MessageBroker interface {
	Publish(ctx context.Context, chatID string, msg []byte) error
}

type UseCase struct {
	messageRepo MessageRepository
	chatRepo    ChatRepository
	broker      MessageBroker
}

func NewUseCase(messageRepo MessageRepository, chatRepo ChatRepository, broker MessageBroker) *UseCase {
	return &UseCase{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		broker:      broker,
	}
}

type SendInput struct {
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Type      domain.MessageType
	Content   string
	ReplyToID *uuid.UUID
}

func (uc *UseCase) Send(ctx context.Context, input SendInput) (*domain.Message, error) {
	isMember, err := uc.chatRepo.IsMember(ctx, input.ChatID, input.SenderID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, domain.ErrForbidden
	}

	if input.Type == domain.MessageTypeText && input.Content == "" {
		return nil, fmt.Errorf("text message content: %w", domain.ErrInvalidInput)
	}

	now := time.Now()
	msg := &domain.Message{
		ID:        uuid.New(),
		ChatID:    &input.ChatID,
		SenderID:  input.SenderID,
		Type:      input.Type,
		Content:   input.Content,
		ReplyToID: input.ReplyToID,
		IsEdited:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return msg, nil
}

func (uc *UseCase) GetHistory(ctx context.Context, chatID, userID uuid.UUID, cursor *pagination.Cursor, limit int) (*pagination.Page[domain.Message], error) {
	isMember, err := uc.chatRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, domain.ErrForbidden
	}

	limit = pagination.NormalizeLimit(limit)

	messages, err := uc.messageRepo.GetChatHistory(ctx, chatID, cursor, limit+1)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	var nextCursor *pagination.Cursor
	if hasMore && len(messages) > 0 {
		last := messages[len(messages)-1]
		nextCursor = &pagination.Cursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID.String(),
		}
	}

	return &pagination.Page[domain.Message]{
		Items:      messages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (uc *UseCase) SearchMessages(ctx context.Context, chatID, userID uuid.UUID, query string, limit int) ([]domain.Message, error) {
	isMember, err := uc.chatRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, domain.ErrForbidden
	}

	limit = pagination.NormalizeLimit(limit)

	messages, err := uc.messageRepo.Search(ctx, chatID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}

	return messages, nil
}
