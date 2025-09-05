package websocket

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

type MessageProvider interface {
	Create(ctx context.Context, message *models.Message) (int64, error)
	UpdateStatus(ctx context.Context, msgId int64, status string)
}

type WebSocketService struct {
	clients        map[string]*Client
	messageService MessageProvider
	mu             sync.RWMutex
}

func NewWebSocketService(messageService MessageProvider) *WebSocketService {
	return &WebSocketService{
		clients:        make(map[string]*Client),
		messageService: messageService,
	}
}

func (s *WebSocketService) HandleConnection(conn *websocket.Conn, userID int64) {
	userIDStr := strconv.FormatInt(userID, 10)
	fmt.Println("userIDStr: ", userIDStr)

	s.mu.Lock()
	if user, exists := s.clients[userIDStr]; exists {
		s.removeClient(user)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte),
		ctx:    ctx,
		cancel: cancel,
	}

	s.clients[userIDStr] = client
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
