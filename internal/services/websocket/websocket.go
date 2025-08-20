package websocket

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/repository"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     int
	UserID int64
	Conn   *websocket.Conn
	Send   chan string
	ctx    context.Context
	cancel context.CancelFunc
}

type WebSocketService struct {
	clients       map[string]*Client
	clientsByConn map[*websocket.Conn]*Client
	userRepo      repository.UserRepository
	mu            sync.Mutex
}

func NewWebSocketService(userRepo repository.UserRepository) *WebSocketService {
	return &WebSocketService{
		clients:       make(map[string]*Client),
		clientsByConn: make(map[*websocket.Conn]*Client),
		userRepo:      userRepo,
	}
}

func (s *WebSocketService) HandleConnection(conn *websocket.Conn, userID int64) {
	userIDStr := strconv.FormatInt(userID, 10)

	s.mu.Lock()
	if user, exists := s.clients[userIDStr]; exists {
		s.removeClient(user.Conn)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		ID:     123, // генерация ID
		UserID: userID,
		Conn:   conn,
		Send:   make(chan string),
		ctx:    ctx,
		cancel: cancel,
	}

	s.clients[userIDStr] = client
	s.clientsByConn[conn] = client
	s.mu.Unlock()

	conn.SetReadLimit(1048576)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go s.handleIncomigMessages(client)
	go s.handleOutgoingMessages(client)

	go s.handlePingPong(client)
}

func (s *WebSocketService) handleIncomigMessages(client *Client) {
	defer func() {
		s.removeClient(client.Conn)
	}()

	for {
		select {
		case <-client.ctx.Done():
			return
		}
	}
}

func (s *WebSocketService) handleOutgoingMessages(client *Client) {
	defer func() {
		s.removeClient(client.Conn)
	}()

	for {
		select {
		case <-client.ctx.Done():
			return
		}
	}
}

func (s *WebSocketService) removeClient(conn *websocket.Conn) {

}

func (s *WebSocketService) handlePingPong(client *Client) {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		s.removeClient(client.Conn)
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
