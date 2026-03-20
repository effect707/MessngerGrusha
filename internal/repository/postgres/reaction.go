package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres/sqlcgen"
)

type ReactionRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewReactionRepository(pool *pgxpool.Pool) *ReactionRepository {
	return &ReactionRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *ReactionRepository) Add(ctx context.Context, reaction *domain.Reaction) error {
	err := r.queries.AddReaction(ctx, sqlcgen.AddReactionParams{
		MessageID: uuidToPgtype(reaction.MessageID),
		UserID:    uuidToPgtype(reaction.UserID),
		Emoji:     reaction.Emoji,
		CreatedAt: timestamptzToPgtype(reaction.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}
	return nil
}

func (r *ReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	err := r.queries.RemoveReaction(ctx, sqlcgen.RemoveReactionParams{
		MessageID: uuidToPgtype(messageID),
		UserID:    uuidToPgtype(userID),
		Emoji:     emoji,
	})
	if err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}
	return nil
}

func (r *ReactionRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]domain.Reaction, error) {
	rows, err := r.queries.GetReactionsByMessageID(ctx, uuidToPgtype(messageID))
	if err != nil {
		return nil, fmt.Errorf("get reactions: %w", err)
	}

	reactions := make([]domain.Reaction, 0, len(rows))
	for _, row := range rows {
		reactions = append(reactions, domain.Reaction{
			MessageID: pgtypeToUUID(row.MessageID),
			UserID:    pgtypeToUUID(row.UserID),
			Emoji:     row.Emoji,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return reactions, nil
}
