package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

type fakeMessageStore struct {
	createCalled   bool
	createdMessage chat.Message
	createErr      error

	getCalled      bool
	requestedRoom  string
	returnMessages []chat.Message
	getErr         error
}

func (fakeStore *fakeMessageStore) CreateMessage(ctx context.Context, message chat.Message) error {
	fakeStore.createCalled = true
	fakeStore.createdMessage = message

	return fakeStore.createErr
}

func (fakeStore *fakeMessageStore) GetMessagesByRoomID(ctx context.Context, roomID string) ([]chat.Message, error) {
	fakeStore.getCalled = true
	fakeStore.requestedRoom = roomID

	if fakeStore.getErr != nil {
		return nil, fakeStore.getErr
	}

	return fakeStore.returnMessages, nil
}

func (fakeStore *fakeMessageStore) GetUserByID(ctx context.Context, id string) (chat.User, error) {
	return chat.User{}, store.ErrUserNotFound
}

func (fakeStore *fakeMessageStore) GetUserByUsername(ctx context.Context, username string) (chat.User, error) {
	return chat.User{}, store.ErrUserNotFound
}

func TestMessageHandlerHistorySuccess(t *testing.T) {
	createdAt := time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC)
	fakeStore := &fakeMessageStore{
		returnMessages: []chat.Message{
			{
				ID:        "message-1",
				UserID:    "user-1",
				RoomID:    "room-123",
				Content:   "hello from history",
				CreatedAt: createdAt,
			},
		},
	}

	handler := newMessageHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/messages?roomId=room-123", nil)
	request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, "user-1"))

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var messageHist messageHistory
	if err := json.NewDecoder(response.Body).Decode(&messageHist); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(messageHist.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messageHist.Messages))
	}

	message := messageHist.Messages[0]
	if message.ID != "message-1" {
		t.Fatalf("expected message ID %q, got %q", "message-1", message.ID)
	}
	if message.RoomID != "room-123" {
		t.Fatalf("expected room ID %q, got %q", "room-123", message.RoomID)
	}
	if message.Content != "hello from history" {
		t.Fatalf("expected message content %q, got %q", "hello from history", message.Content)
	}
	if !message.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created at %v, got %v", createdAt, message.CreatedAt)
	}

	if !fakeStore.getCalled {
		t.Fatal("expected GetMessagesByRoomID to be called")
	}

	if fakeStore.requestedRoom != "room-123" {
		t.Fatalf("expected room ID %q, got %q", "room-123", fakeStore.requestedRoom)
	}
}

func TestMessageHandlerHistoryEmptyRoom(t *testing.T) {
	fakeStore := &fakeMessageStore{}

	handler := newMessageHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/messages?roomId=room-123", nil)
	request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, "user-1"))

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var messageHist messageHistory
	if err := json.NewDecoder(response.Body).Decode(&messageHist); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(messageHist.Messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(messageHist.Messages))
	}
}

func TestMessageHandlerMissingRequiredFields(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
	}{
		{
			name: "missing room id",
			requestBody: `{
			"content": "hi cabin chat"
			}`,
		},
		{
			name: "missing content",
			requestBody: `{
			"roomId": "room-123"}`,
		},
		{
			name:        "missing both fields",
			requestBody: `{}`,
		},
		{
			name: "empty content",
			requestBody: `{
			 "roomId": "room-123",
			 "content": ""
			}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeStore := &fakeMessageStore{}

			handler := newMessageHandler(fakeStore)

			request := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(testCase.requestBody))
			request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, "user-123"))
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

			responseBody := recorder.Body.String()

			if !strings.Contains(responseBody, "roomId and content are required") {
				t.Fatalf("expected response body to contain %q, got %q", "roomId and content are required", responseBody)
			}

			if fakeStore.createCalled {
				t.Fatal("expected CreateMessage to not be called when required fields are missing")
			}

		})
	}

}
