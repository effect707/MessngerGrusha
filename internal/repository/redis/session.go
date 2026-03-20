package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) Set(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	key := sessionKey(userID)
	if err := r.client.Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("set session: %w", err)
	}
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, userID uuid.UUID) (string, error) {
	key := sessionKey(userID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("get session: %w", err)
	}
	return val, nil
}

func (r *SessionRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	key := sessionKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func sessionKey(userID uuid.UUID) string {
	return "session:" + userID.String()
}
