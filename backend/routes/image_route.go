package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/services"
)

func ImageRoute(router fiber.Router, imageService *services.ImageService) {
	router.Get("/images/init", handlers.GetInitialImagesHandler())
	router.Get("/images", handlers.GetImages(imageService))
}
