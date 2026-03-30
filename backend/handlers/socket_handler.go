package handlers

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/sprakash57/snapteil/backend/services"
)

func SocketHandler(socket *services.SocketService) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		socket.Register(conn)
		defer socket.Unregister(conn)

		// Keep connection alive — wait for client to disconnect or send a close frame
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Normal disconnection (close frame, network drop etc.)
				if websocket.IsUnexpectedCloseError(err) {
					log.Printf("WebSocket unexpected close: %v", err)
				}
				break
			}
		}
	}
}
