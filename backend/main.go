package main

import (
	"log"

	"github.com/gofiber/contrib/v3/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/routes"
	"github.com/sprakash57/snapteil/backend/services"
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

	// services
	imageService, err := services.NewImageService("./data/seed.json", "./uploads")
	if err != nil {
		log.Fatal("failed to initialize image service: ", err)
	}

	// routes
	app.Get("/uploads/*", static.New("./uploads"))
	routes.SetupApiV1(app, imageService)

	log.Fatal(app.Listen(":4000"))
}
