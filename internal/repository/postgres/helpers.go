package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func timestamptzToPgtype(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func nullableUUIDToPgtype(id *[16]byte) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}
