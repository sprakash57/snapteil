package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snaptiel/backend/handlers"
)

func main() {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/health", handlers.CheckHealthHandler())

	log.Fatal(app.Listen(":4000"))
}
