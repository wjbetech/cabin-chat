package main

import (
	"encoding/json"
	"net/http"

	"github.com/wjbetech/cabin-chat/pkg/store"
)

type loginHandler struct {
	userStore store.UserStore
	jwtSecret string
}

func newLoginHandler(userStore store.UserStore, jwtSecret string) http.HandlerFunc {
	handler := loginHandler{
		userStore: userStore,
		jwtSecret: jwtSecret,
	}

	return handler.handle
}

func (handler loginHandler) handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var loginReq loginRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&loginReq)

	if err != nil {
		http.Error(writer, "invalid JSON body", http.StatusBadRequest)

		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(writer, "username and/or password required", http.StatusBadRequest)

		return
	}

	http.Error(writer, "not implemented", http.StatusNotImplemented)
}
