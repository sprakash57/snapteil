package handlers

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/services"
)

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
//	@Router			/api/v1/uploads [post]
func UploadImage(imageService *services.ImageService, socket *services.SocketService) fiber.Handler {
	return func(c fiber.Ctx) error {
		title := strings.TrimSpace(c.FormValue("title"))
		if title == "" {
			return fiber.NewError(fiber.StatusBadRequest, "title is required")
		}
		if len([]rune(title)) > 100 {
			return fiber.NewError(fiber.StatusBadRequest, "title must be 100 characters or fewer")
		}

		tags := imageService.ParseTags(c.FormValue("tags"))
		if len(tags) > 5 {
			return fiber.NewError(fiber.StatusBadRequest, "maximum 5 tags allowed")
		}

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

		socket.Broadcast(record)

		log.Printf("Image uploaded: %s (%s)", record.Title, record.Filename)
		return c.Status(fiber.StatusCreated).JSON(record)
	}
}
