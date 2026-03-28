package routes

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/services"
)

func ImageRoute(router fiber.Router, imageService *services.ImageService, hub *services.SocketService) {
	router.Get("/images/init", handlers.GetInitialImagesHandler())
	router.Get("/images", handlers.GetImages(imageService))
	router.Post("/images", handlers.UploadImage(imageService, hub))
}

func WebSocketRoute(router fiber.Router, socket *services.SocketService) {
	router.Get("/ws", websocket.New(handlers.SocketHandler(socket)))
}
