package handlers

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/models"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	// Log the error for debugging purposes
	log.Printf("Error: %v", err)

	// Determine the status code and message based on the error type
	var status int
	var message string

	if e, ok := err.(*fiber.Error); ok {
		status = e.Code
		message = e.Message
	} else {
		status = fiber.StatusInternalServerError
		message = "internal server error"
	}

	// Return a JSON response with the error details
	return c.Status(status).JSON(models.ErrorResponse{
		Status:  status,
		Message: message,
	})
}
