package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/whozdeez/deez-chat2/db"
	"github.com/whozdeez/deez-chat2/handlers"
	"github.com/whozdeez/deez-chat2/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, c *gin.Context) {
	roomID, err := strconv.Atoi(c.Query("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room_id"})
		return
	}
	nickname := c.Query("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nickname is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error: ", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan models.Message, 256),
		roomID:   roomID,
		nickname: nickname,
	}
	hub.register <- client

	messages, err := db.GetMessages(roomID)
	if err != nil {
		log.Println("history error:", err)
	}
	for _, msg := range messages {
		client.send <- msg
	}

	go client.readPump()
	go client.writePump()
}

func main() {
	db.Connect()

	hub := newHub()
	go hub.run()

	r := gin.Default()
	r.Static("/", "./frontend")
	r.GET("/rooms", handlers.GetRooms)
	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server error: ", err)
	}
}