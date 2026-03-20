package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type AttachmentRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewAttachmentRepository(pool *pgxpool.Pool) *AttachmentRepository {
	return &AttachmentRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *AttachmentRepository) Create(ctx context.Context, a *domain.Attachment) error {
	var durationMs pgtype.Int4
	if a.DurationMs != nil {
		durationMs = pgtype.Int4{Int32: int32(*a.DurationMs), Valid: true}
	}

	err := r.queries.CreateAttachment(ctx, sqlcgen.CreateAttachmentParams{
		ID:         uuidToPgtype(a.ID),
		MessageID:  uuidToPgtype(a.MessageID),
		FileName:   a.FileName,
		FileSize:   a.FileSize,
		MimeType:   a.MimeType,
		StorageKey: a.StorageKey,
		DurationMs: durationMs,
		CreatedAt:  timestamptzToPgtype(a.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("create attachment: %w", err)
	}
	return nil
}

func (r *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	row, err := r.queries.GetAttachmentByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get attachment: %w", err)
	}
	return sqlcAttachmentToDomain(row), nil
}

func (r *AttachmentRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Attachment, error) {
	rows, err := r.queries.GetAttachmentsByMessageID(ctx, uuidToPgtype(messageID))
	if err != nil {
		return nil, fmt.Errorf("get attachments: %w", err)
	}

	attachments := make([]domain.Attachment, 0, len(rows))
	for _, row := range rows {
		attachments = append(attachments, *sqlcAttachmentToDomain(row))
	}
	return attachments, nil
}

func sqlcAttachmentToDomain(a sqlcgen.Attachment) *domain.Attachment {
	var durationMs *int
	if a.DurationMs.Valid {
		d := int(a.DurationMs.Int32)
		durationMs = &d
	}

	return &domain.Attachment{
		ID:         pgtypeToUUID(a.ID),
		MessageID:  pgtypeToUUID(a.MessageID),
		FileName:   a.FileName,
		FileSize:   a.FileSize,
		MimeType:   a.MimeType,
		StorageKey: a.StorageKey,
		DurationMs: durationMs,
		CreatedAt:  a.CreatedAt.Time,
	}
}
