package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type Chat struct {
	ID        uuid.UUID `json:"id"`
	Type      ChatType  `json:"type"`
	Name      *string   `json:"name,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

type ChatMember struct {
	ChatID   uuid.UUID  `json:"chat_id"`
	UserID   uuid.UUID  `json:"user_id"`
	Role     MemberRole `json:"role"`
	JoinedAt time.Time  `json:"joined_at"`
}
