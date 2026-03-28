package db

import (
	"context"

	"github.com/whozdeez/deez-chat2/models"
)

func SaveMessage(roomID int, nickname, body string) (models.Message, error) {
	var msg models.Message
	err := Pool.QueryRow(context.Background(),
		`INSERT INTO messages (room_id, nickname, body) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, room_id, nickname, body, created_at`,
		roomID, nickname, body,
	).Scan(&msg.ID, &msg.RoomID, &msg.Nickname, &msg.Body, &msg.CreatedAt)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func GetMessages(roomID int) ([]models.Message, error){
	rows, err := Pool.Query(context.Background(),
		`SELECT id, room_id, nickname, body, created_at
		 FROM messages
		 WHERE room_id = $1
		 ORDER BY created_at ASC
		 LIMIT 50`, 
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.RoomID, &msg.Nickname, &msg.Body, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}