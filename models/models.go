package models

import "time"

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	Nickname  string    `json:"nickname"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
