package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wjbetech/cabin-chat/pkg/store"
)

type sessionHandler struct {
	userStore store.UserStore
}

func (handler sessionHandler) handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	userIDValue := request.Context().Value(authContextKeyUserID)

	if userIDValue == nil {
		http.Error(writer, "missing authenticated user ID", http.StatusUnauthorized)

		return
	}

	userID, authOkay := userIDValue.(string)

	if !authOkay || userID == "" {
		http.Error(writer, "invalid authenticated user ID", http.StatusUnauthorized)

		return
	}

	user, err := handler.userStore.GetUserByID(request.Context(), userID)

	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(writer, "invalid authenticated user - user was not found", http.StatusUnauthorized)

			return
		}

		http.Error(writer, "failed to load user - internal server error", http.StatusInternalServerError)

		return
	}

	response := sessionResponse{
		UserID:     user.ID,
		Username:   user.Username,
		CreatedAt:  user.CreatedAt,
		Status:     user.Status,
		LastSeenAt: user.LastSeenAt,
		AvatarURL:  user.AvatarURL,
	}

	responseBody, err := json.Marshal(response)

	if err != nil {
		http.Error(writer, "failed to json encode response", http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(responseBody)

	if err != nil {
		return
	}
}

func newSessionHandler(userStore store.UserStore) http.HandlerFunc {
	handler := sessionHandler{
		userStore: userStore,
	}

	return handler.handle
}
