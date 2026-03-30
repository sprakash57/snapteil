package routes

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/timeout"
	"github.com/sprakash57/snapteil/backend/handlers"
	"github.com/sprakash57/snapteil/backend/services"
)

func UploadRouteV1(router fiber.Router, imageService *services.ImageService, socket *services.SocketService) {
	router.Post(
		"/uploads",
		limiter.New(limiter.Config{
			Max:        10,
			Expiration: time.Minute,
			LimitReached: func(c fiber.Ctx) error {
				return fiber.NewError(fiber.StatusTooManyRequests, "too many upload attempts")
			},
		}),
		timeout.New(handlers.UploadImage(imageService, socket), timeout.Config{
			Timeout: 30 * time.Second,
			OnTimeout: func(c fiber.Ctx) error {
				return fiber.NewError(fiber.StatusRequestTimeout, "upload timed out")
			},
		}),
	)
}
