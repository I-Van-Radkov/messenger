package http

import (
	"github.com/I-Van-Radkov/messenger/internal/app/http/handlers"
	"github.com/I-Van-Radkov/messenger/internal/app/http/middlwares"
	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouterGin(
	authHandlers *handlers.AuthHandlers,
	wsHandlers *handlers.WebSocketHandlers,
	chatHandlers *handlers.ChatHandlers,
	userHandlers *handlers.UserHandlers,
	cfgAuth *config.AuthConfig,
) *gin.Engine {

	router := gin.Default()

	router.Use(middlwares.CorsMiddleware())

	auth := router.Group("/api")
	{
		auth.POST("/auth/register", authHandlers.RegisterHandler)
		auth.POST("/auth/login", authHandlers.LoginHandler)
		auth.POST("/logout", authHandlers.LogoutHandler)
	}

	protected := router.Group("/api")
	protected.Use(middlwares.JWTAuthMiddleware(cfgAuth.JwtSecret))
	{
		protected.GET("/ws", wsHandlers.WebSocketHandler)

		protected.GET("/users/profile", userHandlers.GetUserHandler)
		protected.GET("/users/search", userHandlers.SearchHandler)

		protected.GET("/chats/", chatHandlers.GetChatsHandler)
		//protected.POST("/chats/:username", chatHandlers.CreateChatHandler)
		protected.GET("/chats/:dialog_id", chatHandlers.GetUserChatHandler)
	}

	return router
}
