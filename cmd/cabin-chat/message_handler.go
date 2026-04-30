package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/id"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

type messageHandler struct {
	messageStore store.MessageStore
}

func newMessageHandler(messageStore store.MessageStore) http.HandlerFunc {
	handler := messageHandler{
		messageStore: messageStore,
	}

	return handler.handle
}

func (handler messageHandler) handle(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		handler.handleCreate(writer, request)
	case http.MethodGet:
		handler.handleHistory(writer, request)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (handler messageHandler) handleCreate(writer http.ResponseWriter, request *http.Request) {
	var messageReq messageRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&messageReq); err != nil {
		http.Error(writer, "invalid JSON body", http.StatusBadRequest)

		return
	}

	if messageReq.RoomID == "" || messageReq.Content == "" {
		http.Error(writer, "roomId and content are required", http.StatusBadRequest)

		return
	}

	userID, err := authenticatedUserID(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)

		return
	}

	message := chat.Message{
		ID:        id.New(),
		UserID:    userID,
		RoomID:    messageReq.RoomID,
		Content:   messageReq.Content,
		CreatedAt: time.Now().UTC(),
	}

	if err := handler.messageStore.CreateMessage(request.Context(), message); err != nil {
		http.Error(writer, "failed to create message", http.StatusInternalServerError)

		return
	}

	response := messageResponse{
		ID:        message.ID,
		UserID:    message.UserID,
		RoomID:    message.RoomID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

func (handler messageHandler) handleHistory(writer http.ResponseWriter, request *http.Request) {
	roomID := request.URL.Query().Get("roomId")
	if roomID == "" {
		http.Error(writer, "roomId is required", http.StatusBadRequest)

		return
	}

	if _, err := authenticatedUserID(request); err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)

		return
	}

	messages, err := handler.messageStore.GetMessagesByRoomID(request.Context(), roomID)
	if err != nil {
		http.Error(writer, "failed to load messages", http.StatusInternalServerError)

		return
	}

	responseMessages := make([]messageResponse, 0, len(messages))
	for _, message := range messages {
		responseMessages = append(responseMessages, messageResponse{
			ID:        message.ID,
			UserID:    message.UserID,
			RoomID:    message.RoomID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		})
	}

	response := messageHistory{
		Messages: responseMessages,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

func authenticatedUserID(request *http.Request) (string, error) {
	userIDValue := request.Context().Value(authContextKeyUserID)
	if userIDValue == nil {
		return "", errors.New("missing authenticated user ID")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return "", errors.New("invalid authenticated user ID")
	}

	return userID, nil
}
