package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type ChatRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *ChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	row, err := r.queries.GetChatByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get chat by id: %w", err)
	}
	return sqlcChatToDomain(row), nil
}

func (r *ChatRepository) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	result, err := r.queries.IsChatMember(ctx, sqlcgen.IsChatMemberParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return false, fmt.Errorf("check membership: %w", err)
	}
	return result, nil
}

func (r *ChatRepository) CreateDirect(ctx context.Context, creatorID, recipientID uuid.UUID) (*domain.Chat, error) {
	chatID := uuid.New()
	now := time.Now()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.queries.WithTx(tx)

	err = q.CreateChat(ctx, sqlcgen.CreateChatParams{
		ID:        uuidToPgtype(chatID),
		Type:      sqlcgen.ChatTypeDirect,
		CreatedBy: uuidToPgtype(creatorID),
		CreatedAt: timestamptzToPgtype(now),
		UpdatedAt: timestamptzToPgtype(now),
	})
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}

	err = q.CreateDirectChatLookup(ctx, sqlcgen.CreateDirectChatLookupParams{
		ChatID: uuidToPgtype(chatID),
		UserA:  uuidToPgtype(creatorID),
		UserB:  uuidToPgtype(recipientID),
	})
	if err != nil {
		return nil, fmt.Errorf("create direct chat lookup: %w", err)
	}

	for _, uid := range []uuid.UUID{creatorID, recipientID} {
		err = q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{
			ChatID:   uuidToPgtype(chatID),
			UserID:   uuidToPgtype(uid),
			Role:     string(domain.MemberRoleMember),
			JoinedAt: timestamptzToPgtype(now),
		})
		if err != nil {
			return nil, fmt.Errorf("add member: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &domain.Chat{
		ID:        chatID,
		Type:      domain.ChatTypeDirect,
		CreatedBy: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *ChatRepository) CreateGroup(ctx context.Context, name string, creatorID uuid.UUID, memberIDs []uuid.UUID) (*domain.Chat, error) {
	chatID := uuid.New()
	now := time.Now()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.queries.WithTx(tx)

	err = q.CreateChat(ctx, sqlcgen.CreateChatParams{
		ID:        uuidToPgtype(chatID),
		Type:      sqlcgen.ChatTypeGroup,
		Name:      textToPgtype(&name),
		CreatedBy: uuidToPgtype(creatorID),
		CreatedAt: timestamptzToPgtype(now),
		UpdatedAt: timestamptzToPgtype(now),
	})
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}

	err = q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{
		ChatID:   uuidToPgtype(chatID),
		UserID:   uuidToPgtype(creatorID),
		Role:     string(domain.MemberRoleAdmin),
		JoinedAt: timestamptzToPgtype(now),
	})
	if err != nil {
		return nil, fmt.Errorf("add creator: %w", err)
	}

	for _, uid := range memberIDs {
		if uid == creatorID {
			continue
		}
		err = q.AddChatMember(ctx, sqlcgen.AddChatMemberParams{
			ChatID:   uuidToPgtype(chatID),
			UserID:   uuidToPgtype(uid),
			Role:     string(domain.MemberRoleMember),
			JoinedAt: timestamptzToPgtype(now),
		})
		if err != nil {
			return nil, fmt.Errorf("add member: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &domain.Chat{
		ID:        chatID,
		Type:      domain.ChatTypeGroup,
		Name:      &name,
		CreatedBy: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *ChatRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	rows, err := r.queries.GetUserChats(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("get user chats: %w", err)
	}

	chats := make([]domain.Chat, 0, len(rows))
	for _, row := range rows {
		chats = append(chats, *sqlcChatToDomain(row))
	}
	return chats, nil
}

func (r *ChatRepository) GetDirectChatID(ctx context.Context, userA, userB uuid.UUID) (uuid.UUID, error) {
	chatID, err := r.queries.GetDirectChat(ctx, sqlcgen.GetDirectChatParams{
		Column1: uuidToPgtype(userA),
		Column2: uuidToPgtype(userB),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, domain.ErrNotFound
		}
		return uuid.UUID{}, fmt.Errorf("get direct chat: %w", err)
	}
	return pgtypeToUUID(chatID), nil
}

func (r *ChatRepository) AddMember(ctx context.Context, chatID, userID uuid.UUID, role domain.MemberRole) error {
	now := time.Now()
	err := r.queries.AddChatMember(ctx, sqlcgen.AddChatMemberParams{
		ChatID:   uuidToPgtype(chatID),
		UserID:   uuidToPgtype(userID),
		Role:     string(role),
		JoinedAt: timestamptzToPgtype(now),
	})
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *ChatRepository) RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error {
	err := r.queries.RemoveChatMember(ctx, sqlcgen.RemoveChatMemberParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (r *ChatRepository) GetMembers(ctx context.Context, chatID uuid.UUID) ([]domain.ChatMember, error) {
	rows, err := r.queries.GetChatMembers(ctx, uuidToPgtype(chatID))
	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}

	members := make([]domain.ChatMember, 0, len(rows))
	for _, row := range rows {
		members = append(members, domain.ChatMember{
			ChatID:   pgtypeToUUID(row.ChatID),
			UserID:   pgtypeToUUID(row.UserID),
			Role:     domain.MemberRole(row.Role),
			JoinedAt: row.JoinedAt.Time,
		})
	}
	return members, nil
}

func sqlcChatToDomain(c sqlcgen.Chat) *domain.Chat {
	var name *string
	if c.Name.Valid {
		name = &c.Name.String
	}
	var avatarURL *string
	if c.AvatarUrl.Valid {
		avatarURL = &c.AvatarUrl.String
	}

	return &domain.Chat{
		ID:        pgtypeToUUID(c.ID),
		Type:      domain.ChatType(c.Type),
		Name:      name,
		AvatarURL: avatarURL,
		CreatedBy: pgtypeToUUID(c.CreatedBy),
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}
