package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wjbetech/cabin-chat/pkg/auth"
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

	user, err := handler.userStore.GetUserByUsername(request.Context(), loginReq.Username)

	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(writer, "invalid username or password", http.StatusUnauthorized)

			return
		}

		http.Error(writer, "failed to load user", http.StatusInternalServerError)

		return
	}

	err = auth.CheckPassword(loginReq.Password, user.HashedPassword)

	if err != nil {
		http.Error(writer, "invalid username or password", http.StatusUnauthorized)

		return
	}

	http.Error(writer, "not implemented", http.StatusNotImplemented)
}
