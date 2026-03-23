package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wjbetech/cabin-chat/pkg/env"
)

func main() {
	config, error := env.Load()
	if error != nil {
		log.Fatalf("load config: %v", error)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := http.StatusOK

		writer.WriteHeader(status)
		_, _ = writer.Write([]byte(
			fmt.Sprintf("status: %d | health: ok | service: cabin-chat", status),
		))
	})

	address := fmt.Sprintf(":%s", config.Port)
	log.Printf("starting cabin-chat server on %s", address)

	if error := http.ListenAndServe(address, mux); error != nil {
		log.Fatalf("server stopped: %v", error)
	}
}
