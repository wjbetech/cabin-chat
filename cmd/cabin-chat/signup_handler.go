package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/wjbetech/cabin-chat/pkg/auth"
	"github.com/wjbetech/cabin-chat/pkg/chat"
	"github.com/wjbetech/cabin-chat/pkg/id"
	"github.com/wjbetech/cabin-chat/pkg/store"
)

type signupHandler struct {
	userStore store.UserStore
	jwtSecret string
}

func newSignupHandler(userStore store.UserStore, jwtSecret string) http.HandlerFunc {
	handler := signupHandler {
		userStore: userStore,
		jwtSecret: jwtSecret,
	}
	
	return handler.handle
}

func (handler signupHandler) handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var signupReq signupRequest
	
	// := declares and inits a new variable
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	
	err := decoder.Decode(&signupReq)
	
	if err != nil {
		http.Error(writer, "invalid JSON body", http.StatusBadRequest)
		
		return
	}
	
	if signupReq.Username == "" || signupReq.Password == "" {
		http.Error(writer, "username and password are required", http.StatusBadRequest)
		
		return
	}
	
	hashedPassword, err := auth.HashPassword(signupReq.Password)
	
	if err != nil {
		http.Error(writer, "failed to hash password", http.StatusInternalServerError)
		
		return
	}
	
	user := chat.User {
		ID: id.New(),
		Username: signupReq.Username,
		HashedPassword: hashedPassword,
		CreatedAt: time.Now().UTC(),
		Status: "offline",
		AvatarURL: "",
	}
	
	err = handler.userStore.CreateUser(request.Context(), user)
	
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyExists) {
			http.Error(writer, "username already exists", http.StatusConflict)
			
			return
		}
		
		http.Error(writer, "failed to create user", http.StatusInternalServerError)
		
		return
	}
	
	token, err := auth.GenerateJWTAccessToken(user.ID, handler.jwtSecret, 24*time.Hour)
	
	if err != nil {
		http.Error(writer, "failed to generate access token", http.StatusInternalServerError)
		
		return
	}
	
	response := signupResponse {
		UserID: user.ID,
		Username: user.Username,
		Token: token,
	}
	
	responseBody, err := json.Marshal(response)
	
	if err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
		
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	
	_, err = writer.Write(responseBody)
	
	if err != nil {
		return
	}
	
}