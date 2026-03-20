package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type ChannelRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewChannelRepository(pool *pgxpool.Pool) *ChannelRepository {
	return &ChannelRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *ChannelRepository) Create(ctx context.Context, ch *domain.Channel) error {
	err := r.queries.CreateChannel(ctx, sqlcgen.CreateChannelParams{
		ID:          uuidToPgtype(ch.ID),
		Slug:        ch.Slug,
		Name:        ch.Name,
		Description: ch.Description,
		AvatarUrl:   textToPgtype(ch.AvatarURL),
		OwnerID:     uuidToPgtype(ch.OwnerID),
		IsPrivate:   ch.IsPrivate,
		CreatedAt:   timestamptzToPgtype(ch.CreatedAt),
		UpdatedAt:   timestamptzToPgtype(ch.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("create channel: %w", err)
	}
	return nil
}

func (r *ChannelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	row, err := r.queries.GetChannelByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get channel: %w", err)
	}
	return sqlcChannelToDomain(row), nil
}

func (r *ChannelRepository) GetBySlug(ctx context.Context, slug string) (*domain.Channel, error) {
	row, err := r.queries.GetChannelBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get channel by slug: %w", err)
	}
	return sqlcChannelToDomain(row), nil
}

func (r *ChannelRepository) Update(ctx context.Context, ch *domain.Channel) error {
	err := r.queries.UpdateChannel(ctx, sqlcgen.UpdateChannelParams{
		ID:          uuidToPgtype(ch.ID),
		Name:        ch.Name,
		Description: ch.Description,
		AvatarUrl:   textToPgtype(ch.AvatarURL),
		IsPrivate:   ch.IsPrivate,
		UpdatedAt:   timestamptzToPgtype(ch.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("update channel: %w", err)
	}
	return nil
}

func (r *ChannelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteChannel(ctx, uuidToPgtype(id)); err != nil {
		return fmt.Errorf("delete channel: %w", err)
	}
	return nil
}

func (r *ChannelRepository) GetPublic(ctx context.Context, limit int) ([]domain.Channel, error) {
	rows, err := r.queries.GetPublicChannels(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("get public channels: %w", err)
	}

	channels := make([]domain.Channel, 0, len(rows))
	for _, row := range rows {
		channels = append(channels, *sqlcChannelToDomain(row))
	}
	return channels, nil
}

func (r *ChannelRepository) GetUserChannels(ctx context.Context, userID uuid.UUID) ([]domain.Channel, error) {
	rows, err := r.queries.GetUserChannels(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user channels: %w", err)
	}

	channels := make([]domain.Channel, 0, len(rows))
	for _, row := range rows {
		channels = append(channels, *sqlcChannelToDomain(row))
	}
	return channels, nil
}

func (r *ChannelRepository) AddSubscriber(ctx context.Context, sub *domain.ChannelSubscriber) error {
	err := r.queries.AddChannelSubscriber(ctx, sqlcgen.AddChannelSubscriberParams{
		ChannelID:    uuidToPgtype(sub.ChannelID),
		UserID:       uuidToPgtype(sub.UserID),
		Role:         string(sub.Role),
		SubscribedAt: timestamptzToPgtype(sub.SubscribedAt),
	})
	if err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}
	return nil
}

func (r *ChannelRepository) RemoveSubscriber(ctx context.Context, channelID, userID uuid.UUID) error {
	err := r.queries.RemoveChannelSubscriber(ctx, sqlcgen.RemoveChannelSubscriberParams{
		ChannelID: uuidToPgtype(channelID),
		UserID:    uuidToPgtype(userID),
	})
	if err != nil {
		return fmt.Errorf("remove subscriber: %w", err)
	}
	return nil
}

func (r *ChannelRepository) IsSubscriber(ctx context.Context, channelID, userID uuid.UUID) (bool, error) {
	result, err := r.queries.IsChannelSubscriber(ctx, sqlcgen.IsChannelSubscriberParams{
		ChannelID: uuidToPgtype(channelID),
		UserID:    uuidToPgtype(userID),
	})
	if err != nil {
		return false, fmt.Errorf("check subscriber: %w", err)
	}
	return result, nil
}

func (r *ChannelRepository) GetSubscribers(ctx context.Context, channelID uuid.UUID) ([]domain.ChannelSubscriber, error) {
	rows, err := r.queries.GetChannelSubscribers(ctx, uuidToPgtype(channelID))
	if err != nil {
		return nil, fmt.Errorf("get subscribers: %w", err)
	}

	subs := make([]domain.ChannelSubscriber, 0, len(rows))
	for _, row := range rows {
		subs = append(subs, domain.ChannelSubscriber{
			ChannelID:    pgtypeToUUID(row.ChannelID),
			UserID:       pgtypeToUUID(row.UserID),
			Role:         domain.SubscriberRole(row.Role),
			SubscribedAt: row.SubscribedAt.Time,
		})
	}
	return subs, nil
}

func sqlcChannelToDomain(c sqlcgen.Channel) *domain.Channel {
	var avatarURL *string
	if c.AvatarUrl.Valid {
		avatarURL = &c.AvatarUrl.String
	}

	return &domain.Channel{
		ID:          pgtypeToUUID(c.ID),
		Slug:        c.Slug,
		Name:        c.Name,
		Description: c.Description,
		AvatarURL:   avatarURL,
		OwnerID:     pgtypeToUUID(c.OwnerID),
		IsPrivate:   c.IsPrivate,
		CreatedAt:   c.CreatedAt.Time,
		UpdatedAt:   c.UpdatedAt.Time,
	}
}
