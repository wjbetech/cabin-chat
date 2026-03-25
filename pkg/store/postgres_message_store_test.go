package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/db"
	"github.com/wjbetech/cabin-chat/pkg/id"
)

func TestPostgresMessageStore_CreateMessage(t *testing.T) {
	databaseURL := os.Getenv("CABIN_CHAT_DATABASE_URL")

	if databaseURL == "" {
		t.Skip("CABIN_CHAT_DATABASE_URL is not set!")
	}

	ctx := context.Background()

	postgresDB, err := db.OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	defer postgresDB.Close()

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, rooms, users CASCADE")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)
	messageStore := NewPostgresMessageStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "message-test-user",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	roomID := id.New()

	_, err = postgresDB.ExecContext(
		ctx,
		`INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`,
		roomID,
		"general",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "hello cabin chat",
		CreatedAt: time.Now().UTC(),
	}

	err = messageStore.CreateMessage(ctx, message)
	if err != nil {
		t.Fatalf("CreateMessage returned an error: %v", err)
	}

	messages, err := messageStore.GetMessagesByRoomID(ctx, roomID)
	if err != nil {
		t.Fatalf("GetMessagesByRoomID returned an error after create: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	savedMessage := messages[0]

	if savedMessage.ID != message.ID {
		t.Fatalf("expected savedMessage.ID to be %s, got %s", message.ID, savedMessage.ID)
	}

	if savedMessage.UserID != message.UserID {
		t.Fatalf("expected savedMessage.UserID to be %s, got %s", message.UserID, savedMessage.UserID)
	}

	if savedMessage.RoomID != message.RoomID {
		t.Fatalf("expected savedMessage.RoomID to be %s, got %s", message.RoomID, savedMessage.RoomID)
	}

	if savedMessage.Content != message.Content {
		t.Fatalf("expected savedMessage.Content to be %s, got %s", message.Content, savedMessage.Content)
	}
}

func TestPostgresMessageStore_GetMessagesByRoomID(t *testing.T) {
	databaseURL := os.Getenv("CABIN_CHAT_DATABASE_URL")

	if databaseURL == "" {
		t.Skip("CABIN_CHAT_DATABASE_URL is not set!")
	}

	ctx := context.Background()

	postgresDB, err := db.OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	defer postgresDB.Close()

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, rooms, users CASCADE")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)
	messageStore := NewPostgresMessageStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "message-room-user",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	roomID := id.New()

	_, err = postgresDB.ExecContext(
		ctx,
		`INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`,
		roomID,
		"general",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	firstMessage := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "first message",
		CreatedAt: time.Now().UTC().Add(-2 * time.Minute),
	}

	secondMessage := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "second message",
		CreatedAt: time.Now().UTC().Add(-1 * time.Minute),
	}

	err = messageStore.CreateMessage(ctx, secondMessage)
	if err != nil {
		t.Fatalf("CreateMessage returned an error for secondMessage: %v", err)
	}

	err = messageStore.CreateMessage(ctx, firstMessage)
	if err != nil {
		t.Fatalf("CreateMessage returned an error for firstMessage: %v", err)
	}

	messages, err := messageStore.GetMessagesByRoomID(ctx, roomID)
	if err != nil {
		t.Fatalf("GetMessagesByRoomID returned an error: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].ID != firstMessage.ID {
		t.Fatalf("expected first result to be %s, got %s", firstMessage.ID, messages[0].ID)
	}

	if messages[1].ID != secondMessage.ID {
		t.Fatalf("expected second result to be %s, got %s", secondMessage.ID, messages[1].ID)
	}
}

func TestPostgresMessageStore_GetMessagesByRoomID_EmptyRoom(t *testing.T) {
	databaseURL := os.Getenv("CABIN_CHAT_DATABASE_URL")

	if databaseURL == "" {
		t.Skip("CABIN_CHAT_DATABASE_URL is not set!")
	}

	ctx := context.Background()

	postgresDB, err := db.OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	defer postgresDB.Close()

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, rooms, users CASCADE")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	messageStore := NewPostgresMessageStore(postgresDB)

	roomID := id.New()

	_, err = postgresDB.ExecContext(
		ctx,
		`INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`,
		roomID,
		"empty-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	messages, err := messageStore.GetMessagesByRoomID(ctx, roomID)
	if err != nil {
		t.Fatalf("GetMessagesByRoomID returned an error for an empty room: %v", err)
	}

	if len(messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(messages))
	}
}

func TestPostgresMessageStore_GetMessagesByRoomID_SameTimestampUsesIDTiebreaker(t *testing.T) {
	databaseURL := os.Getenv("CABIN_CHAT_DATABASE_URL")

	if databaseURL == "" {
		t.Skip("CABIN_CHAT_DATABASE_URL is not set!")
	}

	ctx := context.Background()

	postgresDB, err := db.OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	defer postgresDB.Close()

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, rooms, users CASCADE")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)
	messageStore := NewPostgresMessageStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "message-tiebreak-user",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	roomID := id.New()

	_, err = postgresDB.ExecContext(
		ctx,
		`INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`,
		roomID,
		"tiebreak-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	sameCreatedAt := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)

	firstByID := chat.Message{
		ID:        "11111111-1111-1111-1111-111111111111",
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "first by id",
		CreatedAt: sameCreatedAt,
	}

	secondByID := chat.Message{
		ID:        "22222222-2222-2222-2222-222222222222",
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "second by id",
		CreatedAt: sameCreatedAt,
	}

	err = messageStore.CreateMessage(ctx, secondByID)
	if err != nil {
		t.Fatalf("CreateMessage returned an error for secondByID: %v", err)
	}

	err = messageStore.CreateMessage(ctx, firstByID)
	if err != nil {
		t.Fatalf("CreateMessage returned an error for firstByID: %v", err)
	}

	messages, err := messageStore.GetMessagesByRoomID(ctx, roomID)
	if err != nil {
		t.Fatalf("GetMessagesByRoomID returned an error: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].ID != firstByID.ID {
		t.Fatalf("expected first result to be %s, got %s", firstByID.ID, messages[0].ID)
	}

	if messages[1].ID != secondByID.ID {
		t.Fatalf("expected second result to be %s, got %s", secondByID.ID, messages[1].ID)
	}
}
