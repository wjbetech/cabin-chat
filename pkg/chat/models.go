package chat

import "time"

type User struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	// it's a huge security hole to return the hashed password, so we exclude it from JSON responses
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	Status         string    `json:"status"`
	LastSeenAt     time.Time `json:"lastSeenAt"`
	AvatarURL      string    `json:"avatarUrl"`
}

type Message struct {
	ID        string    `json:"ID"`
	UserID    string    `json:"userId"`
	RoomID    string    `json:"roomId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Reaction struct {
	ID        string    `json:"ID"`
	MessageID string    `json:"messageId"`
	UserID    string    `json:"userId"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"createdAt"`
}
