package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const onlineTTL = 60 * time.Second

type OnlineStatusRepository struct {
	client *redis.Client
}

func NewOnlineStatusRepository(client *redis.Client) *OnlineStatusRepository {
	return &OnlineStatusRepository{client: client}
}

func (r *OnlineStatusRepository) SetOnline(ctx context.Context, userID uuid.UUID) error {
	key := onlineKey(userID)
	if err := r.client.Set(ctx, key, "1", onlineTTL).Err(); err != nil {
		return fmt.Errorf("set online: %w", err)
	}
	return nil
}

func (r *OnlineStatusRepository) IsOnline(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := onlineKey(userID)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check online: %w", err)
	}
	return exists > 0, nil
}

func (r *OnlineStatusRepository) SetOffline(ctx context.Context, userID uuid.UUID) error {
	key := onlineKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("set offline: %w", err)
	}
	return nil
}

func onlineKey(userID uuid.UUID) string {
	return "online:" + userID.String()
}
