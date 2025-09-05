package models

import "time"

type Message struct {
	ID           int64     `json:"id"`
	DialogID     int64     `json:"dialog_id"`
	SenderID     int64     `json:"sender_id"`
	RecipientID  int64     `json:"recipient_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	IsReplyToMsg bool      `json:"is_reply_to_msg"`
	ReplyToMsgID int64     `json:"reply_to_msg_id,omitempty"`
	Status       string    `json:"status"`
}
