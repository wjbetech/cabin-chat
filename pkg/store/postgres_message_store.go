package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
)

type PostgresMessageStore struct {
	db *sql.DB
}

var _ MessageStore = (*PostgresMessageStore)(nil)

func NewPostgresMessageStore(db *sql.DB) *PostgresMessageStore {
	return &PostgresMessageStore{db: db}
}

func (store *PostgresMessageStore) CreateMessage(ctx context.Context, message chat.Message) error {
	if message.ID == "" {
		return fmt.Errorf("message ID is required")
	}

	if message.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	if message.RoomID == "" {
		return fmt.Errorf("room ID is required")
	}

	if message.Content == "" {
		return fmt.Errorf("content is required")
	}

	createdAt := message.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	_, err := store.db.ExecContext(
		ctx,
		`INSERT INTO messages (id, content, user_id, room_id, created_at) VALUES ($1, $2, $3, $4, $5)`,
		message.ID,
		message.Content,
		message.UserID,
		message.RoomID,
		createdAt,
	)

	if err != nil {
		return fmt.Errorf("Create message: %w", err)
	}

	return nil
}

func (store *PostgresMessageStore) GetMessagesByRoomID(ctx context.Context, roomID string) ([]chat.Message, error) {
	rows, err := store.db.QueryContext(
		ctx,
		`SELECT id, content, user_id, room_id, created_at FROM messages WHERE room_id = $1 ORDER BY created_at ASC`,
		roomID,
	)

	if err != nil {
		return nil, fmt.Errorf("query messages by room ID: %w", err)
	}

	defer rows.Close()

	messages := make([]chat.Message, 0)

	for rows.Next() {
		var message chat.Message

		err := rows.Scan(&message.ID, &message.Content, &message.UserID, &message.RoomID, &message.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan message row: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate message rows: %w", err)
	}

	return messages, nil
}
