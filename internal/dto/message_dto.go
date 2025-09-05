package dto

import (
	"time"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type MessageDTO struct {
	ID           int64     `json:"id"`
	SenderID     int64     `json:"sender_id"`
	RecipientID  int64     `json:"recipient_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	IsReplyToMsg bool      `json:"is_reply_to_msg"`
	ReplyToMsgID int64     `json:"reply_to_msg_id,omitempty"`
	Status       string    `json:"status"`
}

type MessagesResponse struct {
	Count    int           `json:"count"`
	Messages []*MessageDTO `json:"messages"`
}

func ToMessageDTO(message *models.Message) *MessageDTO {
	return &MessageDTO{
		ID:           message.ID,
		SenderID:     message.SenderID,
		RecipientID:  message.RecipientID,
		Content:      message.Content,
		CreatedAt:    message.CreatedAt,
		IsReplyToMsg: message.IsReplyToMsg,
		ReplyToMsgID: message.ReplyToMsgID,
		Status:       message.Status,
	}
}
