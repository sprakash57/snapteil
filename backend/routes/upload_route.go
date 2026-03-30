package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/services"
)

func UploadRouteV1(router fiber.Router, imageService *services.ImageService, socket *services.SocketService) {
	router.Post("/uploads", handlers.UploadImage(imageService, socket))
}
