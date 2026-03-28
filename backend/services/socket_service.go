package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

// SocketService manages all active WebSocket connections and broadcasts messages to them.
type SocketService struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

func NewSocketService() *SocketService {
	return &SocketService{
		clients: make(map[*websocket.Conn]bool),
	}
}

// Register adds a new WebSocket client connection
func (socketService *SocketService) Register(conn *websocket.Conn) {
	socketService.mu.Lock()
	defer socketService.mu.Unlock()
	socketService.clients[conn] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(socketService.clients))
}

// Unregister removes a client connection when it disconnects
func (socketService *SocketService) Unregister(conn *websocket.Conn) {
	socketService.mu.Lock()
	defer socketService.mu.Unlock()
	delete(socketService.clients, conn)
	log.Printf("WebSocket client disconnected. Total clients: %d", len(socketService.clients))
}

// Broadcast sends a JSON payload to all connected clients.
// Clients that fail to receive are silently removed.
func (socketService *SocketService) Broadcast(payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("WebSocket broadcast marshal error: %v", err)
		return
	}

	socketService.mu.RLock()
	defer socketService.mu.RUnlock()

	for conn := range socketService.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
	}
}
