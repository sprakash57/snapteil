package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/services"
)

func SetupApiV1(app *fiber.App, imageService *services.ImageService, socket *services.SocketService) {
	apiV1 := app.Group("/api/v1")

	HealthRouteV1(apiV1)
	ImageRouteV1(apiV1, imageService)
	UploadRouteV1(apiV1, imageService, socket)

	// 404 catch-all
	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}

func SetupSocketV1(app *fiber.App, socket *services.SocketService) {
	// WebSocket route for real-time updates
	WebSocketRouteV1(app, socket)
}
