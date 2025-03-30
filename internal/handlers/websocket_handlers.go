package handlers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все источники
	},
}

type WebSocketHandler struct {
	LabServiceURL string // URL для связи с сервисом лабораторий
}

func NewWebSocketHandler(labServiceURL string) *WebSocketHandler {
	return &WebSocketHandler{
		LabServiceURL: labServiceURL,
	}
}
func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем WebSocket соединение
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	// Чтение сообщений от клиента
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Логируем полученное сообщение
		log.Printf("Received message: %s\n", p)

		// Проверяем, что сообщение не пустое
		if len(p) == 0 {
			log.Println("Received empty command")
			conn.WriteMessage(messageType, []byte("Empty command received"))
			continue
		}

		// Преобразуем строку в команду
		command := string(p)

		// Логируем полученную команду
		log.Printf("Parsed command: %s\n", command)

		// Выполнение команды через сервис лабораторий
		labID := 1 // Здесь можно передать labID динамически, если нужно
		output, err := h.ExecuteCommandInLab(labID, command)
		if err != nil {
			conn.WriteMessage(messageType, []byte("Error executing command"))
			continue
		}

		// Отправка результата обратно клиенту
		if err := conn.WriteMessage(messageType, []byte(output)); err != nil {
			log.Println("Error sending response:", err)
			break
		}
	}
}

// Выполнение команды через API сервиса лабораторий
func (h *WebSocketHandler) ExecuteCommandInLab(labID int, command string) (string, error) {
	url := fmt.Sprintf("%s/labs/%d/execute-command", h.LabServiceURL, labID)

	// Создаем тело запроса для команды
	requestBody := fmt.Sprintf(`{"command": "%s"}`, command)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response from lab service: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
