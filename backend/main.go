package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sprakash57/snaptiel/backend/handlers"
)

func main() {
	app := fiber.New()

	// middelewares
	app.Use(logger.New())

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/health", handlers.CheckHealthHandler()) // can be used by load balancers to check if the server is healthy

	log.Fatal(app.Listen(":4000"))
}
