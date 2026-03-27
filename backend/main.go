package main

import (
	"log"

	"github.com/gofiber/contrib/v3/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sprakash57/snaptiel/backend/handlers"
	"github.com/sprakash57/snaptiel/backend/routes"
)

//	@title			Snapteil API
//	@version		1.0
//	@description	Snapteil backend API
//	@host			localhost:4000
//	@BasePath		/

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	// middlewares
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
	}))

	// routes
	routes.SetupApiV1(app)

	log.Fatal(app.Listen(":4000"))
}
