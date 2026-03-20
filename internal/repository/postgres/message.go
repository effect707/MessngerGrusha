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
	"github.com/effect707/MessngerGrusha/internal/pkg/pagination"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type MessageRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *MessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	var chatID, channelID, replyToID pgtype.UUID

	if msg.ChatID != nil {
		chatID = uuidToPgtype(*msg.ChatID)
	}
	if msg.ChannelID != nil {
		channelID = uuidToPgtype(*msg.ChannelID)
	}
	if msg.ReplyToID != nil {
		replyToID = uuidToPgtype(*msg.ReplyToID)
	}

	err := r.queries.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		ID:        uuidToPgtype(msg.ID),
		ChatID:    chatID,
		ChannelID: channelID,
		SenderID:  uuidToPgtype(msg.SenderID),
		Type:      sqlcgen.MessageType(msg.Type),
		Content:   msg.Content,
		ReplyToID: replyToID,
		IsEdited:  msg.IsEdited,
		CreatedAt: timestamptzToPgtype(msg.CreatedAt),
		UpdatedAt: timestamptzToPgtype(msg.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	row, err := r.queries.GetMessageByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get message by id: %w", err)
	}
	return sqlcMessageToDomain(row), nil
}

func (r *MessageRepository) GetChatHistory(ctx context.Context, chatID uuid.UUID, cursor *pagination.Cursor, limit int) ([]domain.Message, error) {
	params := sqlcgen.GetChatHistoryParams{
		ChatID: uuidToPgtype(chatID),
		Limit:  int32(limit),
	}

	if cursor != nil {
		params.CursorCreatedAt = pgtype.Timestamptz{Time: cursor.CreatedAt, Valid: true}
		cursorUUID, err := uuid.Parse(cursor.ID)
		if err != nil {
			return nil, fmt.Errorf("parse cursor id: %w", err)
		}
		params.CursorID = pgtype.UUID{Bytes: cursorUUID, Valid: true}
	}

	rows, err := r.queries.GetChatHistory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get chat history: %w", err)
	}

	messages := make([]domain.Message, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, *historyRowToDomain(row))
	}
	return messages, nil
}

func (r *MessageRepository) Search(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]domain.Message, error) {
	rows, err := r.queries.SearchMessages(ctx, sqlcgen.SearchMessagesParams{
		ChatID:         uuidToPgtype(chatID),
		PlaintoTsquery: query,
		Limit:          int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}

	messages := make([]domain.Message, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, *searchRowToDomain(row))
	}
	return messages, nil
}

func (r *MessageRepository) Update(ctx context.Context, msg *domain.Message) error {
	err := r.queries.UpdateMessage(ctx, sqlcgen.UpdateMessageParams{
		ID:        uuidToPgtype(msg.ID),
		Content:   msg.Content,
		UpdatedAt: timestamptzToPgtype(msg.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("update message: %w", err)
	}
	return nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteMessage(ctx, uuidToPgtype(id))
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	return nil
}

type messageRow struct {
	ID        pgtype.UUID
	ChatID    pgtype.UUID
	ChannelID pgtype.UUID
	SenderID  pgtype.UUID
	Type      sqlcgen.MessageType
	Content   string
	ReplyToID pgtype.UUID
	IsEdited  bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func rowToDomain(m messageRow) *domain.Message {
	msg := &domain.Message{
		ID:       pgtypeToUUID(m.ID),
		SenderID: pgtypeToUUID(m.SenderID),
		Type:     domain.MessageType(m.Type),
		Content:  m.Content,
		IsEdited: m.IsEdited,
	}

	if m.ChatID.Valid {
		id := pgtypeToUUID(m.ChatID)
		msg.ChatID = &id
	}
	if m.ChannelID.Valid {
		id := pgtypeToUUID(m.ChannelID)
		msg.ChannelID = &id
	}
	if m.ReplyToID.Valid {
		id := pgtypeToUUID(m.ReplyToID)
		msg.ReplyToID = &id
	}
	if m.CreatedAt.Valid {
		msg.CreatedAt = m.CreatedAt.Time
	}
	if m.UpdatedAt.Valid {
		msg.UpdatedAt = m.UpdatedAt.Time
	}

	return msg
}

func sqlcMessageToDomain(m sqlcgen.GetMessageByIDRow) *domain.Message {
	return rowToDomain(messageRow{
		ID: m.ID, ChatID: m.ChatID, ChannelID: m.ChannelID, SenderID: m.SenderID,
		Type: m.Type, Content: m.Content, ReplyToID: m.ReplyToID,
		IsEdited: m.IsEdited, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	})
}

func historyRowToDomain(m sqlcgen.GetChatHistoryRow) *domain.Message {
	return rowToDomain(messageRow{
		ID: m.ID, ChatID: m.ChatID, ChannelID: m.ChannelID, SenderID: m.SenderID,
		Type: m.Type, Content: m.Content, ReplyToID: m.ReplyToID,
		IsEdited: m.IsEdited, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	})
}

func searchRowToDomain(m sqlcgen.SearchMessagesRow) *domain.Message {
	return rowToDomain(messageRow{
		ID: m.ID, ChatID: m.ChatID, ChannelID: m.ChannelID, SenderID: m.SenderID,
		Type: m.Type, Content: m.Content, ReplyToID: m.ReplyToID,
		IsEdited: m.IsEdited, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	})
}
