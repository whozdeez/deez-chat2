package db

import (
	"context"
	"log"
)

func Migrate() {
	_, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			room_id INT NOT NULL REFERENCES rooms(id),
			nickname TEXT NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		INSERT INTO rooms (name) VALUES ('General'), ('Random') 
		ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		log.Panic("Migration error: ", err)
	}
}
