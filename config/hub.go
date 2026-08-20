package config

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int
	Conn   *websocket.Conn
	Mu     sync.Mutex
}

type Hub struct {
	Clients map[int]*Client
	Mu      sync.RWMutex
}

func (h *Hub) AddClient(userID int, conn *websocket.Conn) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	h.Clients[userID] = &Client{
		UserID: userID,
		Conn:   conn,
	}
}

func (h *Hub) RemoveClient(userID int) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	delete(h.Clients, userID)
}

func (h *Hub) SendToUser(userID int, event interface{}) error {
	h.Mu.RLock()
	client, exists := h.Clients[userID]
	h.Mu.RUnlock()

	if !exists {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	client.Mu.Lock()
	defer client.Mu.Unlock()

	return client.Conn.WriteMessage(
		websocket.TextMessage,
		data,
	)
}
