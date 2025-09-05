package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/I-Van-Radkov/messenger/internal/dto"
	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/utils"
	"github.com/gin-gonic/gin"
)

type ChatProvider interface {
	GetChats(ctx context.Context, userID int64, limit, offset int) ([]*models.Chat, map[int64]*models.Message, error)
	GetUserChat(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error)
}

type ChatHandlers struct {
	chatService ChatProvider
}

func NewChatHandlers(chatService ChatProvider) *ChatHandlers {
	return &ChatHandlers{
		chatService: chatService,
	}
}

func (h *ChatHandlers) GetChatsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt64("user_id")

	limit, offset, err := utils.GetPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// вызов сервиса и получение данных
	chats, lastMsgInChat, err := h.chatService.GetChats(ctx, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	countChats := len(chats)
	chatsResponse := dto.ChatsResponse{
		Count: countChats,
		Chats: make([]*dto.ChatDTO, countChats),
	}

	for i, chat := range chats {
		lastMessage := lastMsgInChat[chat.ID]

		chatsResponse.Chats[i] = dto.ToChatDTO(chat, lastMessage)
	}

	c.JSON(http.StatusOK, chatsResponse)
}

func (h *ChatHandlers) GetUserChatHandler(c *gin.Context) {
	ctx := c.Request.Context()

	dialogID, err := strconv.ParseInt(c.Param("dialog_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid dialog_id format",
		})
		return
	}

	limit, offset, err := utils.GetPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	messages, err := h.chatService.GetUserChat(ctx, dialogID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	countMessages := len(messages)
	messagesResponse := dto.MessagesResponse{
		Count:    countMessages,
		Messages: make([]*dto.MessageDTO, countMessages),
	}

	for i, message := range messages {
		messagesResponse.Messages[i] = dto.ToMessageDTO(message)
	}

	c.JSON(http.StatusOK, messagesResponse)
}
