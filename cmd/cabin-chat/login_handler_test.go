package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wjbetech/cabin-chat/pkg/auth"
	"github.com/wjbetech/cabin-chat/pkg/chat"
)

func TestLoginHandlerSuccess(t *testing.T) {
	hashedPassword, err := auth.HashPassword("super-secret-test-password")

	if err != nil {
		t.Fatalf("failed to hash password for test setup: %v", err)
	}

	fakeStore := &fakeUserStore{
		returnUser: chat.User{
			ID:             "user12345",
			Username:       "testuser",
			HashedPassword: hashedPassword,
		},
	}

	jwtSecret := "test-secret"

	handler := newLoginHandler(fakeStore, jwtSecret)

	requestBody := `
	 {
		"username": "testuser",
		"password": "super-secret-test-password"
	 }
	`

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var loginResp loginResponse

	err = json.NewDecoder(response.Body).Decode(&loginResp)

	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if loginResp.UserID != "user12345" {
		t.Fatalf("expected user ID %q, got %q", "user12345", loginResp.UserID)
	}

	if loginResp.Username != "testuser" {
		t.Fatalf("expected username %q, got %q", "testuser", loginResp.Username)
	}

	if loginResp.Token == "" {
		t.Fatal("expected a valid token in response")
	}

	if !fakeStore.getUserByUsernameCall {
		t.Fatal("expected GetUserByUsername to be called")
	}
}

func TestLoginHandlerInvalidRequestBody(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
	}{
		{
			name: "malformed json string",
			requestBody: `{
				"username": "testuser",
				"password": "super-secret-test-password",
			}`,
		},
		{
			name: "broken json body",
			requestBody: `{
				"username": "testuser",
				"password": "super-secret-test-password",
				"email": "supersecret@example.com",
			}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeStore := &fakeUserStore{}
			jwtSecret := "test-secret"

			handler := newLoginHandler(fakeStore, jwtSecret)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(testCase.requestBody))
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler(recorder, request)

			response := recorder.Result()

			defer response.Body.Close()

			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

			responseBody := recorder.Body.String()

			if !strings.Contains(responseBody, "invalid JSON body") {
				t.Fatalf("expected response body to contain %q, got %q", "invalid JSON body", responseBody)
			}

			if fakeStore.getUserByUsernameCall {
				t.Fatal("expected GetUserByUsername not to be called when request body is invalid, but it was called")
			}
		},
		)
	}
}
