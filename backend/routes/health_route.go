package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snaptiel/backend/handlers"
)

func HealthRoute(router fiber.Router) {
	router.Get("/health", handlers.CheckHealthHandler())
}
