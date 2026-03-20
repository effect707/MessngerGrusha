package domain

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID         uuid.UUID `json:"id"`
	MessageID  uuid.UUID `json:"message_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	StorageKey string    `json:"storage_key"`
	DurationMs *int      `json:"duration_ms,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
