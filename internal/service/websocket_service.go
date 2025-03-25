package service

import (
	"log/slog"
	"net/http"
	"sync"
	"websocket_service/internal/interfaces"

	"github.com/gorilla/websocket"
	"websocket_service/internal/model"
)

type WebSocketService struct {
	Clients map[string]*model.WebSocketClient
	Mu      sync.Mutex
	Logger  *slog.Logger
}

var _ interfaces.WebSocketServiceInterface = (*WebSocketService)(nil)

func NewWebSocketService(logger *slog.Logger) *WebSocketService {
	return &WebSocketService{
		Clients: make(map[string]*model.WebSocketClient),
		Logger:  logger,
	}
}

func (s *WebSocketService) HandleConnection(w http.ResponseWriter, r *http.Request, userID, labID string) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Разрешаем CORS
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Error("Error while updating WebSocket", slog.Any("error", err))
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	client := &model.WebSocketClient{Conn: conn, LabID: labID}

	s.Mu.Lock()
	s.Clients[userID] = client
	s.Mu.Unlock()

	s.Logger.Info("WebSocket connected", slog.String("userID", userID), slog.String("labID", labID))

	go s.handleMessages(userID)
}

func (s *WebSocketService) handleMessages(userID string) {
	client, exists := s.Clients[userID]
	if !exists {
		return
	}
	defer s.CloseConnection(userID)

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			s.Logger.Warn("Error while reading message", slog.String("userID", userID), slog.Any("error", err))
			break
		}

		s.Logger.Info("Message received", slog.String("userID", userID), slog.String("message", string(msg)))
	}
}

func (s *WebSocketService) BroadcastMessage(labID string, message []byte) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for userID, client := range s.Clients {
		if client.LabID == labID {
			if err := client.SendMessage(websocket.TextMessage, message); err != nil {
				s.Logger.Warn("Error when sending message", slog.String("userID", userID), slog.Any("error", err))
			}
		}
	}
}

func (s *WebSocketService) CloseConnection(userID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if client, exists := s.Clients[userID]; exists {
		client.Close()
		delete(s.Clients, userID)
		s.Logger.Info("WebSocket disconnected", slog.String("userID", userID))
	}
}
