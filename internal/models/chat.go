package models

import "time"

type Chat struct {
	ID            int64     `json:"id"`
	User1ID       int64     `json:"user1_id"`
	User2ID       int64     `json:"user2_id"`
	CreatedAt     time.Time `json:"created_at"`
	LastMessageID int64     `json:"last_message_id"`
}
