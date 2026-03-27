package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func CheckHealthHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	}
}
