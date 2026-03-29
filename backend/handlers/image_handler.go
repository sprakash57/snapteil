package handlers

import (
	"encoding/json"
	"log"
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

// UploadImage godoc
//
//	@Summary		Upload an image
//	@Description	Uploads a new image with title and tags
//	@Tags			images
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			title	formData	string	true	"Image title"
//	@Param			tags	formData	string	false	"Comma-separated tags"
//	@Param			file	formData	file	true	"Image file (JPEG, PNG, GIF, WebP, AVIF, SVG)"
//	@Success		201	{object}	models.Image
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/images [post]
func UploadImage(imageService *services.ImageService, socket *services.SocketService) fiber.Handler {
	return func(c fiber.Ctx) error {
		title := strings.TrimSpace(c.FormValue("title"))
		if title == "" {
			return fiber.NewError(fiber.StatusBadRequest, "title is required")
		}

		tags := imageService.ParseTags(c.FormValue("tags"))

		file, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		if file.Size > imageService.MaxFileSize() {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "file size exceeds limit")
		}

		if !imageService.IsValidImageType(file.Header.Get("Content-Type")) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid file type; accepted: JPEG, PNG, GIF, WebP, AVIF, SVG")
		}

		record, err := imageService.Upload(file, title, tags)
		if err != nil {
			log.Printf("Upload failed: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to process image")
		}

		// Notify all connected WebSocket clients about the new image
		socket.Broadcast(record)

		log.Printf("Image uploaded: %s (%s)", record.Title, record.Filename)
		return c.Status(fiber.StatusCreated).JSON(record)
	}
}
