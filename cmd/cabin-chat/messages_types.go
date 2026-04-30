package main

import "time"

type messageRequest struct {
	RoomID  string `json:"roomId"`
	Content string `json:"content"`
}

type messageResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	RoomID    string    `json:"roomId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type messageHistory struct {
	Messages []messageResponse `json:"messages"`
}
