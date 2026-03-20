package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type ChannelRepository interface {
	Create(ctx context.Context, ch *domain.Channel) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Channel, error)
	Update(ctx context.Context, ch *domain.Channel) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPublic(ctx context.Context, limit int) ([]domain.Channel, error)
	GetUserChannels(ctx context.Context, userID uuid.UUID) ([]domain.Channel, error)
	AddSubscriber(ctx context.Context, sub *domain.ChannelSubscriber) error
	RemoveSubscriber(ctx context.Context, channelID, userID uuid.UUID) error
	IsSubscriber(ctx context.Context, channelID, userID uuid.UUID) (bool, error)
	GetSubscribers(ctx context.Context, channelID uuid.UUID) ([]domain.ChannelSubscriber, error)
}

type UseCase struct {
	channelRepo ChannelRepository
}

func NewUseCase(channelRepo ChannelRepository) *UseCase {
	return &UseCase{channelRepo: channelRepo}
}

type CreateInput struct {
	Slug        string
	Name        string
	Description string
	IsPrivate   bool
	OwnerID     uuid.UUID
}

func (uc *UseCase) Create(ctx context.Context, input CreateInput) (*domain.Channel, error) {
	if input.Slug == "" {
		return nil, fmt.Errorf("slug is required: %w", domain.ErrInvalidInput)
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required: %w", domain.ErrInvalidInput)
	}

	if _, err := uc.channelRepo.GetBySlug(ctx, input.Slug); err == nil {
		return nil, fmt.Errorf("slug already taken: %w", domain.ErrAlreadyExists)
	}

	now := time.Now()
	ch := &domain.Channel{
		ID:          uuid.New(),
		Slug:        input.Slug,
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
		IsPrivate:   input.IsPrivate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.channelRepo.Create(ctx, ch); err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}

	sub := &domain.ChannelSubscriber{
		ChannelID:    ch.ID,
		UserID:       input.OwnerID,
		Role:         domain.SubscriberRoleAdmin,
		SubscribedAt: now,
	}
	if err := uc.channelRepo.AddSubscriber(ctx, sub); err != nil {
		return nil, fmt.Errorf("add owner as subscriber: %w", err)
	}

	return ch, nil
}

func (uc *UseCase) GetByID(ctx context.Context, channelID, userID uuid.UUID) (*domain.Channel, error) {
	ch, err := uc.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	if ch.IsPrivate {
		isSub, err := uc.channelRepo.IsSubscriber(ctx, channelID, userID)
		if err != nil {
			return nil, fmt.Errorf("check subscription: %w", err)
		}
		if !isSub {
			return nil, domain.ErrForbidden
		}
	}

	return ch, nil
}

func (uc *UseCase) GetBySlug(ctx context.Context, slug string, userID uuid.UUID) (*domain.Channel, error) {
	ch, err := uc.channelRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	if ch.IsPrivate {
		isSub, err := uc.channelRepo.IsSubscriber(ctx, ch.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("check subscription: %w", err)
		}
		if !isSub {
			return nil, domain.ErrForbidden
		}
	}

	return ch, nil
}

type UpdateInput struct {
	Name        string
	Description string
	IsPrivate   bool
}

func (uc *UseCase) Update(ctx context.Context, channelID, userID uuid.UUID, input UpdateInput) (*domain.Channel, error) {
	ch, err := uc.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	if ch.OwnerID != userID {
		return nil, domain.ErrForbidden
	}

	ch.Name = input.Name
	ch.Description = input.Description
	ch.IsPrivate = input.IsPrivate
	ch.UpdatedAt = time.Now()

	if err := uc.channelRepo.Update(ctx, ch); err != nil {
		return nil, fmt.Errorf("update channel: %w", err)
	}
	return ch, nil
}

func (uc *UseCase) Delete(ctx context.Context, channelID, userID uuid.UUID) error {
	ch, err := uc.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if ch.OwnerID != userID {
		return domain.ErrForbidden
	}

	if err := uc.channelRepo.Delete(ctx, channelID); err != nil {
		return fmt.Errorf("delete channel: %w", err)
	}
	return nil
}

func (uc *UseCase) Subscribe(ctx context.Context, channelID, userID uuid.UUID) error {
	ch, err := uc.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if ch.IsPrivate {
		return fmt.Errorf("cannot subscribe to private channel: %w", domain.ErrForbidden)
	}

	sub := &domain.ChannelSubscriber{
		ChannelID:    channelID,
		UserID:       userID,
		Role:         domain.SubscriberRoleSubscriber,
		SubscribedAt: time.Now(),
	}

	if err := uc.channelRepo.AddSubscriber(ctx, sub); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	return nil
}

func (uc *UseCase) Unsubscribe(ctx context.Context, channelID, userID uuid.UUID) error {
	ch, err := uc.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}

	if ch.OwnerID == userID {
		return fmt.Errorf("owner cannot unsubscribe: %w", domain.ErrInvalidInput)
	}

	if err := uc.channelRepo.RemoveSubscriber(ctx, channelID, userID); err != nil {
		return fmt.Errorf("unsubscribe: %w", err)
	}
	return nil
}

func (uc *UseCase) GetPublicChannels(ctx context.Context, limit int) ([]domain.Channel, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	channels, err := uc.channelRepo.GetPublic(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get public channels: %w", err)
	}
	return channels, nil
}

func (uc *UseCase) GetUserChannels(ctx context.Context, userID uuid.UUID) ([]domain.Channel, error) {
	channels, err := uc.channelRepo.GetUserChannels(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user channels: %w", err)
	}
	return channels, nil
}
