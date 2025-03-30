package main

import (
	"log"
	"net/http"
	"websocket_service/internal/handlers"
)

func main() {

	labServiceURL := "http://localhost:8083"

	websocketHandler := handlers.NewWebSocketHandler(labServiceURL)

	http.HandleFunc("/ws", websocketHandler.HandleConnections)

	log.Println("WebSocket server started on :8087")
	err := http.ListenAndServe(":8087", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
