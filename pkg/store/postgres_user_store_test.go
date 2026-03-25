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

func TestPostgresUserStore_CreateUser(t *testing.T) {

	databaseURL := os.Getenv("CABIN_CHAT_DATABASE_URL")

	// if DB env is not set, the test skips
	// go test Apps. remains stable
	// - good pattern for integration-style tests
	if databaseURL == "" {
		t.Skip("CABIN_CHAT_DATABASE_URL is not set!")
	}

	ctx := context.Background()

	// reuse real DB bootstrap code from db/postgres.go
	// - no duplicate connection logic
	// - tests against the same DB connection path as the app
	postgresDB, err := db.OpenPostgres(ctx, databaseURL)

	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	defer postgresDB.Close()

	// clean up all tables before testing
	// - keep tests repeatable
	// - avoid conflicts with prior runs
	// - duplicate failures easier to reason about
	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	user := chat.User{
		ID:             id.New(),
		Username:       "testuser",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, user)

	if err != nil {
		// %v prints a standard value
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	savedUser, err := userStore.GetUserByID(ctx, user.ID)

	if err != nil {
		t.Fatalf("GetUserByID returned an error after create: %v", err)
	}

	if savedUser.ID != user.ID {
		// %s prints a standard string value
		t.Fatalf("expected savedUser.ID to be %s, got %s", user.ID, savedUser.ID)
	}

	if savedUser.Username != user.Username {
		t.Fatalf("expected savedUser.Username to be %s, got %s", user.Username, savedUser.Username)
	}

	if savedUser.HashedPassword != user.HashedPassword {
		t.Fatalf("expected savedUser.HashedPassword to be %s, got %s", user.HashedPassword, savedUser.HashedPassword)
	}
}

func TestPostgresUserStore_GetUserByID(t *testing.T) {
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

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	expectedUser := chat.User{
		ID:             id.New(),
		Username:       "testuser",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, expectedUser)

	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	actualUser, err := userStore.GetUserByID(ctx, expectedUser.ID)

	if err != nil {
		t.Fatalf("GetUserByID returned an error after create: %v", err)
	}

	if actualUser.ID != expectedUser.ID {
		t.Fatalf("expected actualUser.ID to be %s, got %s", expectedUser.ID, actualUser.ID)
	}

	if actualUser.Username != expectedUser.Username {
		t.Fatalf("expected actualUser.Username to be %s, got %s", expectedUser.Username, actualUser.Username)
	}

	if actualUser.HashedPassword != expectedUser.HashedPassword {
		t.Fatalf("expected actualUser.HashedPassword to be %s, got %s", expectedUser.HashedPassword, actualUser.HashedPassword)
	}

	if actualUser.Status != expectedUser.Status {
		t.Fatalf("expected actualUser.Status to be %s, got %s", expectedUser.Status, actualUser.Status)
	}

	if actualUser.AvatarURL != expectedUser.AvatarURL {
		t.Fatalf("expected actualUser.AvatarURL to be %s, got %s", expectedUser.AvatarURL, actualUser.AvatarURL)
	}
}

func TestPostgresUserStore_GetUserByUsername(t *testing.T) {
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

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	expectedUser := chat.User{
		ID:             id.New(),
		Username:       "testuser",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, expectedUser)

	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	actualUser, err := userStore.GetUserByUsername(ctx, expectedUser.Username)

	if err != nil {
		t.Fatalf("GetUserByUsername returned an error after create: %v", err)
	}

	if actualUser.ID != expectedUser.ID {
		t.Fatalf("expected actualUser.ID to be %s, got %s", expectedUser.ID, actualUser.ID)
	}

	if actualUser.Username != expectedUser.Username {
		t.Fatalf("expected actualUser.Username to be %s, got %s", expectedUser.Username, actualUser.Username)
	}

	if actualUser.HashedPassword != expectedUser.HashedPassword {
		t.Fatalf("expected actualUser.HashedPassword to be %s, got %s", expectedUser.HashedPassword, actualUser.HashedPassword)
	}

	if actualUser.Status != expectedUser.Status {
		t.Fatalf("expected actualUser.Status to be %s, got %s", expectedUser.Status, actualUser.Status)
	}

	if actualUser.AvatarURL != expectedUser.AvatarURL {
		t.Fatalf("expected actualUser.AvatarURL to be %s, got %s", expectedUser.AvatarURL, actualUser.AvatarURL)
	}
}

func TestPostgresUserStore_CreateUser_DuplicationError(t *testing.T) {
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

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	firstUser := chat.User{
		ID:             id.New(),
		Username:       "testuser",
		HashedPassword: "testpassword",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	secondUser := chat.User{
		ID:             id.New(),
		Username:       "testuser", // same as firstUser
		HashedPassword: "testpassword2",
		CreatedAt:      time.Now().UTC(),
		Status:         "offline",
		AvatarURL:      "",
	}

	err = userStore.CreateUser(ctx, firstUser)

	if err != nil {
		t.Fatalf("CreateUser returned an error on first insert: %v", err)
	}

	err = userStore.CreateUser(ctx, secondUser)

	if errors.Is(err, ErrUserAlreadyExists) == false {
		t.Fatal("expected CreateUser to return an error on duplicate insert, but got nil")
	}
}

func TestPostgresUserStore_MissingUserID(t *testing.T) {
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

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	_, err = userStore.GetUserByID(ctx, id.New())

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatal("expected GetUserByID to return an error when user is not found, but got nil")
	}
}

func TestPostgresUserStore_MissingUsername(t *testing.T) {
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

	_, err = postgresDB.ExecContext(ctx, "TRUNCATE TABLE reactions, messages, users CASCADE")

	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	userStore := NewPostgresUserStore(postgresDB)

	_, err = userStore.GetUserByUsername(ctx, "nonexistentuser")

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatal("expected GetUserByUsername to return an error when user is not found, but got nil")
	}
}
