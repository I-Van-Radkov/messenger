package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketProvider interface {
	HandleConnection(conn *websocket.Conn, userID int64)
}

type WebSocketHandlers struct {
	upgrader         websocket.Upgrader
	websocketService WebSocketProvider
}

func NewWebSocketHandlers(websocketService WebSocketProvider) *WebSocketHandlers {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WebSocketHandlers{
		websocketService: websocketService,
		upgrader:         upgrader,
	}
}

func (h *WebSocketHandlers) WebSocketHandler(c *gin.Context) {
	userId := c.Request.Context().Value("user_id").(int64)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.websocketService.HandleConnection(conn, userId)
}
