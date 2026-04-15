package main

import "time"

type sessionResponse struct {
	UserID     string    `json:"userId"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"`
	LastSeenAt time.Time `json:"lastSeenAt"`
	AvatarURL  string    `json:"avatarUrl"`
}
