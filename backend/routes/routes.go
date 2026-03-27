package routes

import (
	"github.com/gofiber/fiber/v3"
)

func SetupApiV1(app *fiber.App) {
	apiV1 := app.Group("/api/v1")

	HealthRoute(apiV1)
	ImageRoute(apiV1)

	// 404 catch-all
	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}
