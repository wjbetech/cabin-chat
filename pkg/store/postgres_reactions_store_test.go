package store

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/db"
	"github.com/wjbetech/cabin-chat/pkg/id"
)

func TestPostgresReactionStore_AddReaction(t *testing.T) {
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
	reactionStore := NewPostgresReactionStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "reaction-test-user",
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
		"reaction-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "message for reaction",
		CreatedAt: time.Now().UTC(),
	}

	err = messageStore.CreateMessage(ctx, message)
	if err != nil {
		t.Fatalf("CreateMessage returned an error: %v", err)
	}

	reaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "fire",
		CreatedAt: time.Now().UTC(),
	}

	err = reactionStore.AddReaction(ctx, reaction)
	if err != nil {
		t.Fatalf("AddReaction returned an error: %v", err)
	}

	reactions, err := reactionStore.GetReactionsByMessageID(ctx, message.ID)
	if err != nil {
		t.Fatalf("GetReactionsByMessageID returned an error: %v", err)
	}

	if len(reactions) != 1 {
		t.Fatalf("expected 1 reaction, got %d", len(reactions))
	}

	savedReaction := reactions[0]

	if savedReaction.ID != reaction.ID {
		t.Fatalf("expected savedReaction.ID to be %s, got %s", reaction.ID, savedReaction.ID)
	}

	if savedReaction.MessageID != reaction.MessageID {
		t.Fatalf("expected savedReaction.MessageID to be %s, got %s", reaction.MessageID, savedReaction.MessageID)
	}

	if savedReaction.UserID != reaction.UserID {
		t.Fatalf("expected savedReaction.UserID to be %s, got %s", reaction.UserID, savedReaction.UserID)
	}

	if savedReaction.Emoji != reaction.Emoji {
		t.Fatalf("expected savedReaction.Emoji to be %s, got %s", reaction.Emoji, savedReaction.Emoji)
	}
}

func TestPostgresReactionStore_GetReactionsByMessageID(t *testing.T) {
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
	reactionStore := NewPostgresReactionStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "reaction-fetch-user",
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
		"reaction-fetch-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "message for fetching reactions",
		CreatedAt: time.Now().UTC(),
	}

	err = messageStore.CreateMessage(ctx, message)
	if err != nil {
		t.Fatalf("CreateMessage returned an error: %v", err)
	}

	firstReaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "fire",
		CreatedAt: time.Now().UTC().Add(-2 * time.Minute),
	}

	secondReaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "thumbs_up",
		CreatedAt: time.Now().UTC().Add(-1 * time.Minute),
	}

	err = reactionStore.AddReaction(ctx, secondReaction)
	if err != nil {
		t.Fatalf("AddReaction returned an error for secondReaction: %v", err)
	}

	err = reactionStore.AddReaction(ctx, firstReaction)
	if err != nil {
		t.Fatalf("AddReaction returned an error for firstReaction: %v", err)
	}

	reactions, err := reactionStore.GetReactionsByMessageID(ctx, message.ID)
	if err != nil {
		t.Fatalf("GetReactionsByMessageID returned an error: %v", err)
	}

	if len(reactions) != 2 {
		t.Fatalf("expected 2 reactions, got %d", len(reactions))
	}

	if reactions[0].ID != firstReaction.ID {
		t.Fatalf("expected first result to be %s, got %s", firstReaction.ID, reactions[0].ID)
	}

	if reactions[1].ID != secondReaction.ID {
		t.Fatalf("expected second result to be %s, got %s", secondReaction.ID, reactions[1].ID)
	}
}

func TestPostgresReactionStore_AddReaction_DuplicateRejected(t *testing.T) {
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
	reactionStore := NewPostgresReactionStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "reaction-duplicate-user",
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
		"reaction-duplicate-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "message for duplicate reaction test",
		CreatedAt: time.Now().UTC(),
	}

	err = messageStore.CreateMessage(ctx, message)
	if err != nil {
		t.Fatalf("CreateMessage returned an error: %v", err)
	}

	firstReaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "fire",
		CreatedAt: time.Now().UTC(),
	}

	secondReaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "fire",
		CreatedAt: time.Now().UTC(),
	}

	err = reactionStore.AddReaction(ctx, firstReaction)
	if err != nil {
		t.Fatalf("AddReaction returned an error for firstReaction: %v", err)
	}

	err = reactionStore.AddReaction(ctx, secondReaction)
	if !errors.Is(err, ErrUserAlreadyReactedWithEmoji) {
		t.Fatalf("expected ErrUserAlreadyReactedWithEmoji, got %v", err)
	}
}

func TestPostgresReactionStore_RemoveReaction(t *testing.T) {
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
	reactionStore := NewPostgresReactionStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "reaction-remove-user",
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
		"reaction-remove-room",
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("insert room: %v", err)
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    user.ID,
		RoomID:    roomID,
		Content:   "message for reaction removal",
		CreatedAt: time.Now().UTC(),
	}

	err = messageStore.CreateMessage(ctx, message)
	if err != nil {
		t.Fatalf("CreateMessage returned an error: %v", err)
	}

	reaction := chat.Reaction{
		ID:        id.New(),
		MessageID: message.ID,
		UserID:    user.ID,
		Emoji:     "fire",
		CreatedAt: time.Now().UTC(),
	}

	err = reactionStore.AddReaction(ctx, reaction)
	if err != nil {
		t.Fatalf("AddReaction returned an error: %v", err)
	}

	err = reactionStore.RemoveReaction(ctx, reaction.ID)
	if err != nil {
		t.Fatalf("RemoveReaction returned an error: %v", err)
	}

	reactions, err := reactionStore.GetReactionsByMessageID(ctx, message.ID)
	if err != nil {
		t.Fatalf("GetReactionsByMessageID returned an error after remove: %v", err)
	}

	if len(reactions) != 0 {
		t.Fatalf("expected 0 reactions after removal, got %d", len(reactions))
	}
}
