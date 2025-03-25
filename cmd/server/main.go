package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"websocket_service/internal/handlers"
	"websocket_service/internal/routes"
	"websocket_service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	webSocketService := service.NewWebSocketService(logger)

	webSocketHandler := handlers.NewWebSocketHandler(webSocketService)

	router := gin.Default()

	routes.SetupWebSocketRoutes(router, webSocketHandler)

	port := "8085"
	logger.Info("WebSocket server is running", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("Error while starting WebSocket", slog.Any("error", err))
	}
}
