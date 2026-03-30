package handlers

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/models"
	"github.com/sprakash57/snapteil/backend/services"
)

// GetInitialImagesHandler godoc
//
//	@Summary		Get initial images
//	@Description	Returns the list of initial images with their metadata
//	@Tags			images
//	@Produce		json
//	@Success		200	{array}		models.Image
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/images/init [get]
func GetInitialImagesHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		data, err := os.ReadFile("./data/seed.json") // inital set of images
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

// GetImages godoc
//
//	@Summary		Get paginated images
//	@Description	Returns a paginated list of images with optional tag filtering
//	@Tags			images
//	@Produce		json
//	@Param			page	query	int	false	"Page number (default: 1)"
//	@Param			perPage	query	int	false	"Items per page (default: 5, max: 20)"
//	@Param			tags	query	string	false	"Comma-separated tags to filter by"
//	@Success		200	{object}	models.PaginatedResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/images [get]
func GetImages(imageService *services.ImageService) fiber.Handler {
	return func(c fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		perPage, _ := strconv.Atoi(c.Query("perPage", "5"))
		tagsParam := c.Query("tags", "")

		var tags []string
		if tagsParam != "" {
			for _, t := range strings.Split(tagsParam, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					tags = append(tags, t)
				}
			}
		}

		if page < 1 || perPage < 1 || perPage > 20 {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid pagination parameters")
		}

		resp := imageService.GetPaginated(page, perPage, tags)
		return c.JSON(resp)
	}
}
