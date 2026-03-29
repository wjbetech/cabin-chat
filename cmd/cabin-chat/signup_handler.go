package main

import (
	"encoding/json"
	"net/http"

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
	
	http.Error(writer, "not implemented", http.StatusNotImplemented)
}