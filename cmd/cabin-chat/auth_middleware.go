package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/wjbetech/cabin-chat/pkg/auth"
)

type authContextKey string

const authContextKeyUserID authContextKey = "authenticatedUserID"

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

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ParseAndValidateJWTAccessToken(token, middleware.jwtSecret)

		if err != nil {
			http.Error(writer, "malformed, expired or invalid token", http.StatusUnauthorized)

			return
		}

		userID := claims.Subject

		requestContext := context.WithValue(request.Context(), authContextKeyUserID, userID)

		next.ServeHTTP(writer, request.WithContext(requestContext))
	})
}
