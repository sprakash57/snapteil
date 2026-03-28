package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/services"
)

func SetupApiV1(app *fiber.App, imageService *services.ImageService) {
	apiV1 := app.Group("/api/v1")

	HealthRoute(apiV1)
	ImageRoute(apiV1, imageService)

	// 404 catch-all
	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}
