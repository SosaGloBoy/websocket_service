package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn  *websocket.Conn
	LabID string
	Mu    sync.Mutex
}

func (c *WebSocketClient) SendMessage(messageType int, message []byte) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return c.Conn.WriteMessage(messageType, message)
}

func (c *WebSocketClient) Close() {
	c.Conn.Close()
}
