package pagination

import (
	"time"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

type Page[T any] struct {
	Items      []T     `json:"items"`
	NextCursor *Cursor `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}

const DefaultLimit = 50
const MaxLimit = 100

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}
