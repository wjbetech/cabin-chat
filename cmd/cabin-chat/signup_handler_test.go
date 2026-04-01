package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wjbetech/cabin-chat/pkg/auth"
	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

type fakeUserStore struct {
	createdUser chat.User
	createErr error
	createCalled bool
}

func (fakeStore *fakeUserStore) CreateUser(ctx context.Context, user chat.User) error {
	fakeStore.createCalled = true
	fakeStore.createdUser = user

	return fakeStore.createErr	
}

func (fakeStore *fakeUserStore) GetUserByID(ctx context.Context, id string) (chat.User, error) {
	return chat.User{}, store.ErrUserNotFound
}

func (fakeStore *fakeUserStore) GetUserByUsername(ctx context.Context, username string) (chat.User, error) {
	return chat.User{}, store.ErrUserNotFound
}

func TestSignupHandlerSuccess(t *testing.T) {
	fakeStore := &fakeUserStore{}
	jwtSecret := "test-secret"
	
	handler := newSignupHandler(fakeStore, jwtSecret)
	
	requestBody := `
	{
		"username": "testuser",
		"password": "super-secret-test-password"
	}`
	
	request := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	
	handler(recorder, request)
	
	response := recorder.Result()
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
	}
	
	var signupResp signupResponse
	
	err := json.NewDecoder(response.Body).Decode(&signupResp)
	
	
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	
	if signupResp.UserID == "" {
		t.Fatal("expected non-empty user ID in response, got empty string")
	}
	
	if signupResp.Username != "testuser" {
		t.Fatalf("expected username 'testuser' in response, got '%s'", signupResp.Username)
	}
	
	if signupResp.Token == "" {
		t.Fatal("expected non-empty token in response, got empty string")
	}
	
	if fakeStore.createdUser.ID == "" {
		t.Fatal("expected created user to have non-empty ID, got empty string")
	}
	
	if fakeStore.createdUser.Username != "testuser" {
		t.Fatalf("expected created user to have username 'testuser', got '%s'", fakeStore.createdUser.Username)
	}
	
	if fakeStore.createdUser.HashedPassword == "" {
		t.Fatal("expected created user to have non-empty hashed password, got empty string")
	}
	
	if fakeStore.createdUser.HashedPassword == "super-secret-test-password" {
		t.Fatal("expected created user's hashed password to not be the same as the plaintext password")
	}
	
	err = auth.CheckPassword("super-secret-test-password", fakeStore.createdUser.HashedPassword)
	
	if err != nil {
		t.Fatalf("expected password to match hashed password, got error: %v", err)
	}
}

func TestSignupHandlerInvalidRequestBody(t *testing.T) {
	testCases := []struct {
		name string
		requestBody string
	} {
		{
		name: "malformed json",
		requestBody:
		`
			{
				"username": "testuser",
				"password": "super-secret-test-password",
			}
		`,
	},
	{
		name: "missing required fields",
		requestBody:
		`
			{
				"username": "testuser"
				"password": "super-secret-test-password"
				"email": "supersecret@example.com
			}
		`,
	},
}
	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			
			fakeStore := &fakeUserStore{}
			jwtSecret := "test-secret"
			
			handler := newSignupHandler(fakeStore, jwtSecret)
			
			request := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(testCase.requestBody))
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
				t.Fatalf("expected response body to contain 'invalid JSON body', got '%s'", responseBody)
			}
			
			if fakeStore.createCalled {
				t.Fatal("expected CreateUser not to be called when request body is invalid, but it was called")
			}
		},
	)}
}

func TestSignupHandlerMissingRequiredFields(t *testing.T) {
	testCases := []struct {
		name string
		requestBody string
	} {
		{
			name: "missing username",
			requestBody: `
			 {
				"password": "super-secret-test-password"
			 }
			`,
		},
		{
			name: "missing password",
			requestBody: `
			 {
				"username": "testuser"
			 }
			`,
		},
		{
			name: "missing both username and password",
			requestBody: `{}`,
		},
		{
			name: "empty username",
			requestBody: `
			 {
				"username": "",
				"password": "super-secret-test-password"
			 }
			`,
		},
		{
			name: "empty password",
			requestBody: `
			 {
				"username": "testuser",
				"password": ""
			 }
			`,
		},
	}
	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeStore := &fakeUserStore{}
			jwtSecret := "test-secret"
			
			handler := newSignupHandler(fakeStore, jwtSecret)
			
			request := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(testCase.requestBody))
			
			request.Header.Set("Content-Type", "application/json")
			
			recorder := httptest.NewRecorder()
			
			handler(recorder, request)
			
			response := recorder.Result()
			
			defer response.Body.Close()
			
			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}
			
			responseBody := recorder.Body.String()
			
			if !strings.Contains(responseBody, "username and password are required") {
				t.Fatalf("expected response body to contain 'username and password are required', got '%s'", responseBody)
			}
			
			if fakeStore.createCalled {
				t.Fatal("expected CreateUser not to be called when required fields are missing, but it was called")
			}
		})
	}
}

func TestSignupHandlerDuplicateUsername(t *testing.T) {
    testCases := []struct {
        name      string
        createErr error
    }{
        {
            name:      "direct error",
            createErr: store.ErrUserAlreadyExists,
        },
        {
            name:      "wrapped error",
            createErr: fmt.Errorf("wrapped create failure: %w", store.ErrUserAlreadyExists),
        },
    }

    for _, testCase := range testCases {
        t.Run(testCase.name, func(t *testing.T) {
            fakeStore := &fakeUserStore{
                createErr: testCase.createErr,
            }
            jwtSecret := "test-secret"

            handler := newSignupHandler(fakeStore, jwtSecret)

            requestBody := `
            {
                "username": "testuser",
                "password": "super-secret-test-password"
            }`

            request := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
            request.Header.Set("Content-Type", "application/json")

            recorder := httptest.NewRecorder()

            handler(recorder, request)

            response := recorder.Result()
            defer response.Body.Close()

            if response.StatusCode != http.StatusConflict {
                t.Fatalf("expected status %d, got %d", http.StatusConflict, response.StatusCode)
            }

            responseBody := recorder.Body.String()

            if !strings.Contains(responseBody, "username already exists") {
                t.Fatalf("expected response body to contain %q, got %q", "username already exists", responseBody)
            }

            if !fakeStore.createCalled {
                t.Fatal("expected CreateUser to be called for a valid signup request")
            }
        })
    }
}