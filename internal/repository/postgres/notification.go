package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type NotificationRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	err := r.queries.CreateNotification(ctx, sqlcgen.CreateNotificationParams{
		ID:        uuidToPgtype(n.ID),
		UserID:    uuidToPgtype(n.UserID),
		Type:      string(n.Type),
		Payload:   n.Payload,
		IsRead:    n.IsRead,
		CreatedAt: timestamptzToPgtype(n.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	rows, err := r.queries.GetUserNotifications(ctx, sqlcgen.GetUserNotificationsParams{
		UserID: uuidToPgtype(userID),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}

	return sqlcNotificationsToDomain(rows), nil
}

func (r *NotificationRepository) GetUnread(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	rows, err := r.queries.GetUnreadNotifications(ctx, sqlcgen.GetUnreadNotificationsParams{
		UserID: uuidToPgtype(userID),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get unread notifications: %w", err)
	}

	return sqlcNotificationsToDomain(rows), nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	err := r.queries.MarkNotificationRead(ctx, sqlcgen.MarkNotificationReadParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return fmt.Errorf("mark read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.MarkAllNotificationsRead(ctx, uuidToPgtype(userID)); err != nil {
		return fmt.Errorf("mark all read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUnreadNotifications(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, fmt.Errorf("count unread: %w", err)
	}
	return count, nil
}

func sqlcNotificationsToDomain(rows []sqlcgen.Notification) []domain.Notification {
	notifications := make([]domain.Notification, 0, len(rows))
	for _, row := range rows {
		notifications = append(notifications, domain.Notification{
			ID:        pgtypeToUUID(row.ID),
			UserID:    pgtypeToUUID(row.UserID),
			Type:      domain.NotificationType(row.Type),
			Payload:   row.Payload,
			IsRead:    row.IsRead,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return notifications
}
