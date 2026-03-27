package routes

import (
	"github.com/gofiber/fiber/v3"
)

func SetupV1(app *fiber.App) {
	apiV1 := app.Group("/api/v1")

	HealthRoute(apiV1)
}
