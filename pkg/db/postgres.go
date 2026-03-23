package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// blank import the postgresql driver and toolkit
	// - this registers the pgx driver
	// - doesn't require any direct call
	_ "github.com/jackc/pgx/v5/stdlib"
)

// this helper func:
// - validates the URL
// - opens the connection
// - configures the connection pool
// - pings the DB
// - returns the DB handle
func OpenPostgres(ctx context.Context, databaseURL string) (*sql.DB, error) {

	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	// don't overthink these values for now
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	// for now just set a one hour max connection lifetime
	db.SetConnMaxLifetime(60 * time.Minute)

	// WithTimeout() helper protects ping from hanging
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres connection: %w", err)
	}

	return db, nil
}
