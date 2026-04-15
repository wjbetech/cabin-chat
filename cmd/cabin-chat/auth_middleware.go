package main

import (
	"net/http"
	"strings"
)

type authMiddleware struct {
	jwtSecret string
}

func newAuthMiddlware(jwtSecret string) authMiddleware {
	return authMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (middleware authMiddleware) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(writer, "missing Authorization header", http.StatusUnauthorized)

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(writer, "invalid Authorization header format", http.StatusUnauthorized)

			return
		}

		http.Error(writer, "not implemented", http.StatusNotImplemented)
	})
}
