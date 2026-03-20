package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeVoice  MessageType = "voice"
	MessageTypeSystem MessageType = "system"
)

type Message struct {
	ID        uuid.UUID   `json:"id"`
	ChatID    *uuid.UUID  `json:"chat_id,omitempty"`
	ChannelID *uuid.UUID  `json:"channel_id,omitempty"`
	SenderID  uuid.UUID   `json:"sender_id"`
	Type      MessageType `json:"type"`
	Content   string      `json:"content"`
	ReplyToID *uuid.UUID  `json:"reply_to_id,omitempty"`
	IsEdited  bool        `json:"is_edited"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
