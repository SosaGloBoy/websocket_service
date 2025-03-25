package interfaces

import (
	"net/http"
)

type WebSocketServiceInterface interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, userID, labID string)
	BroadcastMessage(labID string, message []byte)
	CloseConnection(userID string)
}
