package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/handlers"
)

func HealthRouteV1(router fiber.Router) {
	router.Get("/health", handlers.CheckHealthHandler())
}
