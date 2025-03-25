package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"websocket_service/internal/service"
)

type WebSocketHandler struct {
	Service *service.WebSocketService
}

// NewWebSocketHandler создает новый обработчик WebSocket
func NewWebSocketHandler(service *service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{Service: service}
}

// HandleWebSocket устанавливает WebSocket-соединение
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID := c.GetString("userID") // Берем из middleware API Gateway
	labID := c.Query("lab_id")

	if userID == "" || labID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userID or labID"})
		return
	}

	h.Service.HandleConnection(c.Writer, c.Request, userID, labID)
}
