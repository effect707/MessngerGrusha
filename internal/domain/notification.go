package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeNewMessage    NotificationType = "new_message"
	NotificationTypeNewSubscriber NotificationType = "new_subscriber"
	NotificationTypeChannelPost   NotificationType = "channel_post"
	NotificationTypeReaction      NotificationType = "reaction"
)

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	Type      NotificationType `json:"type"`
	Payload   json.RawMessage  `json:"payload"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}
