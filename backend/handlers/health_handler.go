package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	_ "github.com/sprakash57/snapteil/backend/models"
)

type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Timestamp string `json:"timestamp" example:"2026-03-27T10:30:00Z"`
}

// CheckHealthHandler godoc
//
//	@Summary		Health check
//	@Description	Returns server health status with a timestamp
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/health [get]
func CheckHealthHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	}
}
