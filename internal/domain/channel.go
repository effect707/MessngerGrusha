package domain

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID          uuid.UUID `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	OwnerID     uuid.UUID `json:"owner_id"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubscriberRole string

const (
	SubscriberRoleAdmin      SubscriberRole = "admin"
	SubscriberRoleSubscriber SubscriberRole = "subscriber"
)

type ChannelSubscriber struct {
	ChannelID    uuid.UUID      `json:"channel_id"`
	UserID       uuid.UUID      `json:"user_id"`
	Role         SubscriberRole `json:"role"`
	SubscribedAt time.Time      `json:"subscribed_at"`
}
