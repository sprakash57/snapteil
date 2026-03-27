package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// CheckHealthHandler godoc
//
//	@Summary		Health check
//	@Description	Returns server health status with a timestamp
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/api/v1/health [get]
func CheckHealthHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	}
}
