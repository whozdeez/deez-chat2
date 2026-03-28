package main

import "github.com/whozdeez/deez-chat2/models"

type Hub struct {
	rooms      map[int]map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[int]map[*Client]bool),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
		case client := <-h.unregister:
			if room, ok := h.rooms[client.roomID]; ok {
				delete(room, client)
				close(client.send)
			}
		
		case msg := <-h.broadcast:
			for client := range h.rooms[msg.RoomID] {
				client.send <- msg
			}
		}
	}
}