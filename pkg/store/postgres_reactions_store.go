package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wjbetech/cabin-chat/pkg/chat"
)

type PostgresReactionStore struct {
	db *sql.DB
}

var _ ReactionStore = (*PostgresReactionStore)(nil)

func NewPostgresReactionStore(db *sql.DB) *PostgresReactionStore {
	return &PostgresReactionStore{db: db}
}

func (store *PostgresReactionStore) AddReaction(ctx context.Context, reaction chat.Reaction) error {
	if reaction.ID == "" {
		return fmt.Errorf("reaction ID is required")
	}

	if reaction.MessageID == "" {
		return fmt.Errorf("reaction message ID is required")
	}

	if reaction.UserID == "" {
		return fmt.Errorf("reaction user ID is required")
	}

	if reaction.Emoji == "" {
		return fmt.Errorf("reaction emoji is required")
	}

	createdAt := reaction.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	_, err := store.db.ExecContext(
		ctx,
		`INSERT INTO reactions (id, message_id, user_id, emoji, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
		reaction.ID,
		reaction.MessageID,
		reaction.UserID,
		reaction.Emoji,
		createdAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyReactedWithEmoji
		}

		return fmt.Errorf("add reaction: %w", err)
	}

	return nil
}

func (store *PostgresReactionStore) GetReactionsByMessageID(ctx context.Context, messageID string) ([]chat.Reaction, error) {
	rows, err := store.db.QueryContext(
		ctx,
		`SELECT id, message_id, user_id, emoji, created_at
        FROM reactions
        WHERE message_id = $1
        ORDER BY created_at ASC, id ASC`,
		messageID,
	)
	if err != nil {
		return nil, fmt.Errorf("get reactions by message ID: %w", err)
	}

	defer rows.Close()

	reactions := make([]chat.Reaction, 0)

	for rows.Next() {
		var reaction chat.Reaction

		err := rows.Scan(
			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Emoji,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan reaction row: %w", err)
		}

		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reaction rows: %w", err)
	}

	return reactions, nil
}

func (store *PostgresReactionStore) RemoveReaction(ctx context.Context, reactionID string) error {
	if reactionID == "" {
		return fmt.Errorf("reaction ID is required")
	}

	_, err := store.db.ExecContext(
		ctx,
		`DELETE FROM reactions WHERE id = $1`,
		reactionID,
	)
	if err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}

	return nil
}
