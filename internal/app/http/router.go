package http

import (
	"github.com/I-Van-Radkov/messenger/internal/app/http/handlers"
	"github.com/I-Van-Radkov/messenger/internal/app/http/middlwares"
	"github.com/gin-gonic/gin"
)

func NewRouterGin(authHandlers *handlers.AuthHandlers) *gin.Engine {
	router := gin.Default()

	router.Use(middlwares.CorsMiddleware())

	router.POST("/api/register/", authHandlers.RegisterHandler)
	router.POST("/api/login/", authHandlers.LoginHandler)

	// router.GET("/api/chats/")
	// router.POST("/api/chats/:username")

	return router
}
