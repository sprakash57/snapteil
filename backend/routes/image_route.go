package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snaptiel/backend/handlers"
	"github.com/sprakash57/snaptiel/backend/services"
)

func ImageRoute(router fiber.Router, imageService *services.ImageService) {
	router.Get("/images/init", handlers.GetInitialImagesHandler())
	router.Get("/images", handlers.GetImages(imageService))
}
