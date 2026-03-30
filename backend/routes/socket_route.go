package routes

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/services"
)

func WebSocketRouteV1(router fiber.Router, socket *services.SocketService) {
	router.Get("/ws", websocket.New(handlers.SocketHandler(socket)))
}
