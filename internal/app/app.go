package app

import (
	apphttp "github.com/I-Van-Radkov/messenger/internal/app/http"
	"github.com/I-Van-Radkov/messenger/internal/app/http/handlers"
	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/I-Van-Radkov/messenger/internal/repository/mariadb"
	"github.com/I-Van-Radkov/messenger/internal/services/auth"
	"github.com/I-Van-Radkov/messenger/internal/services/chat"
	"github.com/I-Van-Radkov/messenger/internal/services/message"
	"github.com/I-Van-Radkov/messenger/internal/services/user"
	"github.com/I-Van-Radkov/messenger/internal/services/websocket"
)

type App struct {
	HTTPServer *apphttp.Server
}

func NewApp(cfg *config.Config) *App {
	db := mariadb.MustConnect(cfg.DB)

	userRepo := mariadb.NewUserRepository(db)
	messageRepo := mariadb.NewMessageRepository(db)
	chatRepo := mariadb.NewChatRepo(db)

	userService := user.NewUserService(userRepo)
	messageService := message.NewMessageService(messageRepo)

	chatService := chat.NewChatService(chatRepo, messageService)

	authService := auth.NewService(
		userService,
		cfg.Auth,
	)

	websocketService := websocket.NewWebSocketService(messageService)

	authHttpHandlers := handlers.NewAuthHandlers(authService)
	websocketHandlers := handlers.NewWebSocketHandlers(websocketService)
	chatHandlers := handlers.NewChatHandlers(chatService)
	userHandlers := handlers.NewUserHandlers(userService)

	httpRouter := apphttp.NewRouterGin(
		authHttpHandlers,
		websocketHandlers,
		chatHandlers,
		userHandlers,
		cfg.Auth,
	)

	httpServer, err := apphttp.NewServer(cfg.HTTP.Port, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout, httpRouter)
	if err != nil {
		panic(err)
	}

	return &App{
		HTTPServer: httpServer,
	}
}
