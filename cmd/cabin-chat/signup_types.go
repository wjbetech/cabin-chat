package main

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signupResponse struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
