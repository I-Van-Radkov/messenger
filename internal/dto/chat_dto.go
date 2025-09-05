package dto

import (
	"time"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type ChatDTO struct {
	ID          int64
	User1ID     int64     `json:"user1_id"`
	User2ID     int64     `json:"user2_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastMessage *MessageDTO
}

type ChatsResponse struct {
	Count int        `json:"count"`
	Chats []*ChatDTO `json:"chats"`
}

func ToChatDTO(chat *models.Chat, lastMessage *models.Message) *ChatDTO {
	return &ChatDTO{
		ID:          chat.ID,
		User1ID:     chat.User1ID,
		User2ID:     chat.User2ID,
		CreatedAt:   chat.CreatedAt,
		LastMessage: ToMessageDTO(lastMessage),
	}
}
