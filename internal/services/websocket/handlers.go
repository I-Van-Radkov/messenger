package websocket

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/gorilla/websocket"
)

func (s *WebSocketService) handleIncomigMessages(client *Client) {
	defer func() {
		s.removeClient(client)
	}()

	for {
		var msg IncomingMessage

		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			return
		}

		if msg.Text == "" {
			s.sendError(client, "Send empty message", "empty_msg")
		}

		s.handleMessage(client, msg)
	}
}

func (s *WebSocketService) handleOutgoingMessages(client *Client) {
	defer func() {
		s.removeClient(client)
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				//
				return
			}
		case <-client.ctx.Done():
			return
		}
	}
}

func (s *WebSocketService) handleMessage(client *Client, message IncomingMessage) {
	switch message.Action {
	case actionSendMessage:
		s.handleSendMessage(client, message)
	default:
		// Unknown Action
	}
}

func (s *WebSocketService) handleSendMessage(client *Client, incMessage IncomingMessage) {
	ctx := context.Background()

	message := &models.Message{
		DialogID:     incMessage.DialogID,
		SenderID:     client.UserID,
		RecipientID:  incMessage.RecipientID,
		Content:      incMessage.Text,
		CreatedAt:    time.Now(),
		IsReplyToMsg: incMessage.IsReplyToMsg,
		ReplyToMsgID: incMessage.ReplyToMsgID,
		Status:       statusSent,
	}

	msgId, err := s.messageService.Create(ctx, message)
	if err != nil {
		s.sendError(client, "Failed to save message", "db_error")
		return
	}

	outMessage := OutgoingMessage{
		Action:    actionNewMessage,
		MessageID: msgId,
		Status:    statusSent,
	}

	outMessageBytes, _ := json.Marshal(outMessage)

	delivered := s.sendToUser(client.UserID, outMessageBytes)

	status := statusSent
	if delivered {
		status = statusDelivered

		s.messageService.UpdateStatus(ctx, msgId, status)
	}

	confirmation := OutgoingMessage{
		Action:    actionMessageStatus,
		MessageID: msgId,
		Status:    status,
	}

	confirmationBytes, _ := json.Marshal(confirmation)
	s.sendToClient(client, confirmationBytes)
}

func (s *WebSocketService) sendToUser(userId int64, message []byte) bool {

	if client, exists := s.getClient(userId); exists {
		select {
		case client.Send <- message:
			return true
		default:
			return false
		}
	}

	return false
}

func (s *WebSocketService) sendToClient(client *Client, message []byte) {
	select {
	case client.Send <- message:
	default:
		// Client channel is full
	}
}

func (s *WebSocketService) sendError(client *Client, message, code string) {
	msgError := ErrorMessage{
		Action:  actionError,
		Message: message,
		Code:    code,
	}

	msgErrorBytes, _ := json.Marshal(msgError)
	s.sendToClient(client, msgErrorBytes)
}

func (s *WebSocketService) removeClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userIdStr := strconv.FormatInt(client.UserID, 10)
	if client, exists := s.clients[userIdStr]; exists {
		client.Conn.Close()
		close(client.Send)
		client.cancel()
		delete(s.clients, userIdStr)
	}
}

func (s *WebSocketService) getClient(userID int64) (*Client, bool) {
	userIdStr := strconv.FormatInt(userID, 10)

	s.mu.RLock()
	defer s.mu.RUnlock()
	client, exists := s.clients[userIdStr]

	return client, exists
}

func (s *WebSocketService) handlePingPong(client *Client) {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		s.removeClient(client)
	}()

	for {
		select {
		case <-ticker.C:
			err := client.Conn.WriteControl(websocket.PingMessage, []byte(time.Now().String()), time.Now().Add(10*time.Second))
			if err != nil {
				return
			}
		case <-client.ctx.Done():
			return
		}
	}
}
