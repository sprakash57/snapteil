package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snaptiel/backend/handlers"
)

func ImageRoute(router fiber.Router) {
	router.Get("/images", handlers.GetInitialImagesHandler())
}
