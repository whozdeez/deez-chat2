package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whozdeez/deez-chat2/db"
)

func GetRooms(c *gin.Context) {
	rows, err := db.Pool.Query(c.Request.Context(),
		`SELECT id, name FROM rooms ORDER BY id ASC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch rooms"})
		return
	}
	defer rows.Close()

	var rooms []gin.H
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		rooms = append(rooms, gin.H{"id": id, "name": name})
	}
	c.JSON(http.StatusOK, rooms)
}
