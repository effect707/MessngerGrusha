package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error)
	GetUnread(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}

type UseCase struct {
	notifRepo NotificationRepository
}

func NewUseCase(notifRepo NotificationRepository) *UseCase {
	return &UseCase{notifRepo: notifRepo}
}

func (uc *UseCase) Send(ctx context.Context, userID uuid.UUID, notifType domain.NotificationType, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	n := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Payload:   data,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := uc.notifRepo.Create(ctx, n); err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (uc *UseCase) GetNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	notifications, err := uc.notifRepo.GetByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}
	return notifications, nil
}

func (uc *UseCase) GetUnread(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	notifications, err := uc.notifRepo.GetUnread(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get unread: %w", err)
	}
	return notifications, nil
}

func (uc *UseCase) MarkRead(ctx context.Context, notifID, userID uuid.UUID) error {
	if err := uc.notifRepo.MarkRead(ctx, notifID, userID); err != nil {
		return fmt.Errorf("mark read: %w", err)
	}
	return nil
}

func (uc *UseCase) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := uc.notifRepo.MarkAllRead(ctx, userID); err != nil {
		return fmt.Errorf("mark all read: %w", err)
	}
	return nil
}

func (uc *UseCase) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := uc.notifRepo.CountUnread(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("count unread: %w", err)
	}
	return count, nil
}
