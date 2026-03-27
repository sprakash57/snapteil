package main

import (
	"log"

	"github.com/gofiber/contrib/v3/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sprakash57/snaptiel/backend/routes"
)

//	@title			Snapteil API
//	@version		1.0
//	@description	Snapteil backend API
//	@host			localhost:4000
//	@BasePath		/

func main() {
	app := fiber.New()

	// middlewares
	app.Use(logger.New())

	// swagger
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
	}))

	// routes
	routes.SetupV1(app)

	log.Fatal(app.Listen(":4000"))
}
