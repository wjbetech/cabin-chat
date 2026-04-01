package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wjbetech/cabin-chat/pkg/db"
	"github.com/wjbetech/cabin-chat/pkg/env"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

func main() {
	cfg, err := env.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	databaseStatus := "disabled"

	var userStore store.UserStore

	if cfg.DatabaseURL != "" {
		postgresDB, err := db.OpenPostgres(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("open postgres connection: %v", err)
		}

		defer func() {
			if err := postgresDB.Close(); err != nil {
				log.Printf("close postgres connection: %v", err)
			}
		}()

		userStore = store.NewPostgresUserStore(postgresDB)

		databaseStatus = "connected"
		log.Printf("postgres connection established")
	} else {
		log.Printf("no database URL configured; starting without postgres")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := http.StatusOK

		writer.WriteHeader(status)
		_, _ = fmt.Fprintf(writer,
			"status: %d | health: ok | service: cabin-chat | database: %s", status, databaseStatus)
	})

	if userStore != nil {
		mux.HandleFunc("/signup", newSignupHandler(userStore, cfg.JWTSecret))
	}

	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("starting cabin-chat server on %s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
