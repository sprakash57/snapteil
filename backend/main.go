package main

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/contrib/v3/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sprakash57/snapteil/backend/config"
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
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
		BodyLimit:    int(cfg.MaxFileSize),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// middlewares
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       300,
	}))
	if cfg.EnableSwagger {
		app.Use(swagger.New(swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
		}))
	}

	// services
	imageService, err := services.NewImageService(cfg)
	if err != nil {
		log.Fatal("failed to initialize image service: ", err)
	}
	socketService := services.NewSocketService()

	// routes
	app.Get("/uploads/*", static.New("./uploads", static.Config{
		ModifyResponse: func(c fiber.Ctx) error {
			c.Set("X-Content-Type-Options", "nosniff")
			if strings.EqualFold(filepath.Ext(c.Path()), ".svg") {
				c.Set("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox")
			}
			return nil
		},
	}))
	routes.SetupSocketV1(app, socketService, cfg)
	routes.SetupApiV1(app, imageService, socketService, cfg)

	log.Fatal(app.Listen(":" + cfg.Port))
}
