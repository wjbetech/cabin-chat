package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

type fakeSessionStore struct {
	getUserByIDCalled bool
	returnUser        chat.User
	getUserByIDErr    error
}

func (fakeStore *fakeSessionStore) CreateUser(ctx context.Context, user chat.User) error {
	return nil
}

func (fakeStore *fakeSessionStore) GetUserByID(ctx context.Context, id string) (chat.User, error) {
	fakeStore.getUserByIDCalled = true

	if fakeStore.getUserByIDErr != nil {
		return chat.User{}, fakeStore.getUserByIDErr
	}

	return fakeStore.returnUser,
		nil
}

func (fakeStore *fakeSessionStore) GetUserByUsername(ctx context.Context, username string) (chat.User, error) {
	return chat.User{}, store.ErrUserNotFound
}

func TestSessionHandlerMissingContextUserID(t *testing.T) {
	fakeStore := &fakeSessionStore{}

	handler := newSessionHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/session", nil)

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
	}

	responseBody := recorder.Body.String()

	if !strings.Contains(responseBody, "missing authenticated user ID") {
		t.Fatalf("expected response body to contain %q, got %q", "missing authenticated user ID", responseBody)
	}

	if fakeStore.getUserByIDCalled {
		t.Fatal("expected GetUserByID to not get called when the user ID is missing from the context")
	}
}

func TestSessionHandlerInvalidContextUserID(t *testing.T) {
	fakeStore := &fakeSessionStore{}

	handler := newSessionHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/session", nil)

	request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, 12345))

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
	}

	responseBody := recorder.Body.String()

	if !strings.Contains(responseBody, "invalid authenticated user ID") {
		t.Fatalf("expected response body to contain %q, got %q", "invalid authenticated user ID", responseBody)
	}

	if fakeStore.getUserByIDCalled {
		t.Fatal("expected GetUserByID not to be called when the context user ID has the wrong type")
	}
}

func TestSessionHandlerEmptyContextUserID(t *testing.T) {
	fakeStore := &fakeSessionStore{}

	handler := newSessionHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/session", nil)

	request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, ""))

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
	}

	responseBody := recorder.Body.String()

	if !strings.Contains(responseBody, "invalid authenticated user ID") {
		t.Fatalf("expected response body to contain %q, got %q", "invalid authenticated user ID", responseBody)
	}

	if fakeStore.getUserByIDCalled {
		t.Fatal("expected GetUserByID not to be called when the context user ID is empty")
	}
}

func TestSessionHandlerMissingUserReturnedFromStore(t *testing.T) {
	testCases := []struct {
		name           string
		getUserByIDErr error
	}{
		{
			name:           "direct not found error",
			getUserByIDErr: store.ErrUserNotFound,
		},
		{
			name:           "wrapped not found error",
			getUserByIDErr: fmt.Errorf("wrapped lookup failure: %w", store.ErrUserNotFound),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeStore := &fakeSessionStore{
				getUserByIDErr: testCase.getUserByIDErr,
			}

			handler := newSessionHandler(fakeStore)

			request := httptest.NewRequest(http.MethodGet, "/session", nil)

			request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, "user-123"))

			recorder := httptest.NewRecorder()

			handler(recorder, request)

			response := recorder.Result()

			defer response.Body.Close()

			if response.StatusCode != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
			}

			responseBody := recorder.Body.String()

			if !strings.Contains(responseBody, "invalid authenticated user - user was not found") {
				t.Fatalf("expected response body to contain %q, got %q", "invalid authenticated user - user was not found", responseBody)
			}

			if !fakeStore.getUserByIDCalled {
				t.Fatal("expected GetUserByID to be called when the context user ID is present")
			}
		})
	}
}

func TestSessionHandlerSuccess(t *testing.T) {
	fakeStore := &fakeSessionStore{
		returnUser: chat.User{
			ID:         "user-123",
			Username:   "testuser",
			CreatedAt:  time.Date(2026, 4, 29, 12, 0, 0, 0, time.UTC),
			Status:     "online",
			LastSeenAt: time.Date(2026, 4, 29, 12, 30, 0, 0, time.UTC),
			AvatarURL:  "https://example.com/avatar.png",
		},
	}

	handler := newSessionHandler(fakeStore)

	request := httptest.NewRequest(http.MethodGet, "/session", nil)

	request = request.WithContext(context.WithValue(request.Context(), authContextKeyUserID, "user-123"))

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var sessionResp sessionResponse
	err := json.NewDecoder(response.Body).Decode(&sessionResp)

	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if sessionResp.UserID != "user-123" {
		t.Fatalf("expected user ID %q, got %q", "user-123", sessionResp.UserID)
	}

	if sessionResp.Username != "testuser" {
		t.Fatalf("expected username %q, got %q", "testuser", sessionResp.Username)
	}

	if !sessionResp.CreatedAt.Equal(time.Date(2026, 4, 29, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected created at %v, got %v", time.Date(2026, 4, 29, 12, 0, 0, 0, time.UTC), sessionResp.CreatedAt)
	}

	if sessionResp.Status != "online" {
		t.Fatalf("expected status %q, got %q", "online", sessionResp.Status)
	}

	if !sessionResp.LastSeenAt.Equal(time.Date(2026, 4, 29, 12, 30, 0, 0, time.UTC)) {
		t.Fatalf("expected last seen at %v, got %v", time.Date(2026, 4, 29, 12, 30, 0, 0, time.UTC), sessionResp.LastSeenAt)
	}

	if sessionResp.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("expected avatar URL %q, got %q", "https://example.com/avatar.png", sessionResp.AvatarURL)
	}

	if !fakeStore.getUserByIDCalled {
		t.Fatal("expected GetUserByID to be called when the context user ID is present")
	}
}
