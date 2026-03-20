package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const typingTTL = 5 * time.Second

type TypingRepository struct {
	client *redis.Client
}

func NewTypingRepository(client *redis.Client) *TypingRepository {
	return &TypingRepository{client: client}
}

func (r *TypingRepository) SetTyping(ctx context.Context, chatID, userID uuid.UUID) error {
	key := typingKey(chatID)
	pipe := r.client.Pipeline()
	pipe.SAdd(ctx, key, userID.String())
	pipe.Expire(ctx, key, typingTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("set typing: %w", err)
	}
	return nil
}

func (r *TypingRepository) GetTyping(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	key := typingKey(chatID)
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get typing: %w", err)
	}

	userIDs := make([]uuid.UUID, 0, len(members))
	for _, m := range members {
		id, err := uuid.Parse(m)
		if err != nil {
			continue
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

func (r *TypingRepository) StopTyping(ctx context.Context, chatID, userID uuid.UUID) error {
	key := typingKey(chatID)
	if err := r.client.SRem(ctx, key, userID.String()).Err(); err != nil {
		return fmt.Errorf("stop typing: %w", err)
	}
	return nil
}

func typingKey(chatID uuid.UUID) string {
	return "typing:" + chatID.String()
}
