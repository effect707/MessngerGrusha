package reaction

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type ReactionRepository interface {
	Add(ctx context.Context, reaction *domain.Reaction) error
	Remove(ctx context.Context, messageID, userID uuid.UUID, emoji string) error
	GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Reaction, error)
}

type MessageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
}

type ChatRepository interface {
	IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

type UseCase struct {
	reactionRepo ReactionRepository
	msgRepo      MessageRepository
	chatRepo     ChatRepository
}

func NewUseCase(reactionRepo ReactionRepository, msgRepo MessageRepository, chatRepo ChatRepository) *UseCase {
	return &UseCase{
		reactionRepo: reactionRepo,
		msgRepo:      msgRepo,
		chatRepo:     chatRepo,
	}
}

func (uc *UseCase) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	if emoji == "" {
		return fmt.Errorf("emoji is required: %w", domain.ErrInvalidInput)
	}

	msg, err := uc.msgRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("get message: %w", err)
	}

	if msg.ChatID != nil {
		isMember, err := uc.chatRepo.IsMember(ctx, *msg.ChatID, userID)
		if err != nil {
			return fmt.Errorf("check membership: %w", err)
		}
		if !isMember {
			return domain.ErrForbidden
		}
	}

	reaction := &domain.Reaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}

	if err := uc.reactionRepo.Add(ctx, reaction); err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}
	return nil
}

func (uc *UseCase) RemoveReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	if err := uc.reactionRepo.Remove(ctx, messageID, userID, emoji); err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}
	return nil
}

func (uc *UseCase) GetReactions(ctx context.Context, messageID, userID uuid.UUID) ([]domain.Reaction, error) {
	msg, err := uc.msgRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}

	if msg.ChatID != nil {
		isMember, err := uc.chatRepo.IsMember(ctx, *msg.ChatID, userID)
		if err != nil {
			return nil, fmt.Errorf("check membership: %w", err)
		}
		if !isMember {
			return nil, domain.ErrForbidden
		}
	}

	reactions, err := uc.reactionRepo.GetByMessageID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get reactions: %w", err)
	}
	return reactions, nil
}
