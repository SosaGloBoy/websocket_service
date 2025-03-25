package routes

import (
	"github.com/gin-gonic/gin"
	"websocket_service/internal/handlers"
)

func SetupWebSocketRoutes(router *gin.Engine, handler *handlers.WebSocketHandler) {
	router.GET("/ws", handler.HandleWebSocket)
}
