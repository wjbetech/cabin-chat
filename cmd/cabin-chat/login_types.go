package main

// matches the two fields the handler needs to auth
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// mirrors signup's safe response shape
type loginResponse struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// no need for HashedPassword here because login responses should never expose password material, even hashed
