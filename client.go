package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/whozdeez/deez-chat2/db"
	"github.com/whozdeez/deez-chat2/models"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan models.Message
	roomID   int
	nickname string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg models.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		saved, err := db.SaveMessage(c.roomID, c.nickname, msg.Body)
		if err != nil {
			log.Println("Error saving message: ", err)
			continue
		}
		c.hub.broadcast <- saved
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error writing message: ", err)
			break
		}
	}
}
