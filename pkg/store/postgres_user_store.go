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

// this struct holds the database handle, meaning:
// - the store depends on the DB connection
// - methods can use store.db to run SQL
type PostgresUserStore struct {
	db *sql.DB
}

// this is a compile-time interface check, a useful
// trick to use in Go, meaning:
//   - if PostgresUserStore ever stops matching UserStore,
//     the compiler will complain immediately
var _ UserStore = (*PostgresUserStore)(nil)

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (store *PostgresUserStore) CreateUser(ctx context.Context, user chat.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}

	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	if user.HashedPassword == "" {
		return fmt.Errorf("hashed password is required")
	}

	createdAt := user.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	status := user.Status
	if status == "" {
		status = "offline"
	}

	var lastSeenAt any
	if user.LastSeenAt.IsZero() {
		lastSeenAt = nil
	} else {
		lastSeenAt = user.LastSeenAt
	}

	_, err := store.db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, hashed_password, status, last_seen_at, avatar_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID,
		user.Username,
		user.HashedPassword,
		status,
		lastSeenAt,
		user.AvatarURL,
		createdAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (store *PostgresUserStore) GetUserByID(ctx context.Context, id string) (chat.User, error) {
	row := store.db.QueryRowContext(
		ctx,
		`SELECT id, username, hashed_password, status, last_seen_at, avatar_url, created_at
		FROM users
		WHERE id = $1`,
		id,
	)
	var user chat.User

	// the schema allows for last_seen_at to be NULL, but
	// the Go model uses LastSeenAt time.Time.
	// plain time.Time cannot represent SQL NULL,
	// so we use sql.NullTime to handle that
	var lastSeenAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Status,
		&lastSeenAt,
		&user.AvatarURL,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return chat.User{},
			ErrUserNotFound
	}

	if err != nil {
		return chat.User{}, fmt.Errorf("get user by ID: %w", err)
	}

	if lastSeenAt.Valid {
		user.LastSeenAt = lastSeenAt.Time
	}

	return user, nil
}

func (store *PostgresUserStore) GetUserByUsername(ctx context.Context, username string) (chat.User, error) {
	row := store.db.QueryRowContext(
		ctx,
		`SELECT id, username, hashed_password, status, last_seen_at, avatar_url, created_at
		FROM users
		WHERE username = $1`,
		username,
	)

	var user chat.User
	var lastSeenAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Status,
		&lastSeenAt,
		&user.AvatarURL,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return chat.User{},
			ErrUserNotFound
	}

	if err != nil {
		return chat.User{}, fmt.Errorf("get user by username: %w", err)
	}

	if lastSeenAt.Valid {
		user.LastSeenAt = lastSeenAt.Time
	}

	return user, nil
}
