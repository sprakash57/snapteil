package handlers

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snaptiel/backend/models"
)

// GetInitialImagesHandler godoc
//
//	@Summary		Get initial images
//	@Description	Returns the list of initial images with their metadata
//	@Tags			images
//	@Produce		json
//	@Success		200	{array}		models.Image
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/images [get]
func GetInitialImagesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		data, err := os.ReadFile("./data/seed.json")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load images")
		}

		var images []models.Image
		if err := json.Unmarshal(data, &images); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to parse image data")
		}

		return c.JSON(images)
	}
}
