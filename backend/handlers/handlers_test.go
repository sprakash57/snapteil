package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sprakash57/snapteil/backend/config"
	"github.com/sprakash57/snapteil/backend/models"
	"github.com/sprakash57/snapteil/backend/services"
)

func TestErrorHandlerReturnsFiberErrorPayload(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/fiber", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	})

	resp := performRequest(t, app, httptest.NewRequest("GET", "/fiber", nil))

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var body models.ErrorResponse
	decodeJSON(t, resp, &body)

	if body.Status != fiber.StatusBadRequest || body.Message != "bad request" {
		t.Fatalf("body = %#v, want status=400 message=bad request", body)
	}
}

func TestErrorHandlerReturnsGenericInternalServerError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/generic", func(c fiber.Ctx) error {
		return errors.New("boom")
	})

	resp := performRequest(t, app, httptest.NewRequest("GET", "/generic", nil))

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusInternalServerError)
	}

	var body models.ErrorResponse
	decodeJSON(t, resp, &body)

	if body.Status != fiber.StatusInternalServerError || body.Message != "internal server error" {
		t.Fatalf("body = %#v, want status=500 message=internal server error", body)
	}
}

func TestCheckHealthHandlerReturnsOkWithTimestamp(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/health", CheckHealthHandler())

	resp := performRequest(t, app, httptest.NewRequest("GET", "/health", nil))

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body HealthResponse
	decodeJSON(t, resp, &body)

	if body.Status != "ok" {
		t.Fatalf("Status = %q, want ok", body.Status)
	}
	if _, err := time.Parse(time.RFC3339, body.Timestamp); err != nil {
		t.Fatalf("Timestamp = %q, want RFC3339 timestamp: %v", body.Timestamp, err)
	}
}

func TestGetInitialImagesHandlerReturnsSeedImages(t *testing.T) {
	restore := chdirToBackendRoot(t)
	defer restore()

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/images/init", GetInitialImagesHandler())

	resp := performRequest(t, app, httptest.NewRequest("GET", "/images/init", nil))

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var images []models.Image
	decodeJSON(t, resp, &images)

	if len(images) == 0 {
		t.Fatal("expected seeded images to be returned")
	}
	if images[0].Filename == "" {
		t.Fatal("expected seeded image filename to be populated")
	}
}

func TestGetImagesRejectsInvalidPagination(t *testing.T) {
	restore := chdirToBackendRoot(t)
	defer restore()

	imageService := mustNewImageService(t)
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/images", GetImages(imageService))

	resp := performRequest(t, app, httptest.NewRequest("GET", "/images?perPage=21", nil))

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var body models.ErrorResponse
	decodeJSON(t, resp, &body)

	if body.Message != "Invalid pagination parameters" {
		t.Fatalf("Message = %q, want Invalid pagination parameters", body.Message)
	}
}

func TestUploadImageRejectsInvalidImageFile(t *testing.T) {
	restore := chdirToBackendRoot(t)
	defer restore()

	imageService := mustNewImageService(t)
	socketService := services.NewSocketService()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("title", "Invalid upload"); err != nil {
		t.Fatalf("WriteField(title) error = %v", err)
	}
	if err := writer.WriteField("tags", "test"); err != nil {
		t.Fatalf("WriteField(tags) error = %v", err)
	}
	part, err := writer.CreateFormFile("file", "notes.txt")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write([]byte("plain text is not an image")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/uploads", UploadImage(imageService, socketService))

	req := httptest.NewRequest("POST", "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := performRequest(t, app, req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var errResp models.ErrorResponse
	decodeJSON(t, resp, &errResp)

	if errResp.Message != "invalid or unsupported image file" {
		t.Fatalf("Message = %q, want invalid or unsupported image file", errResp.Message)
	}
}

func TestUploadImageRejectsTagWithInvalidCharacters(t *testing.T) {
	restore := chdirToBackendRoot(t)
	defer restore()

	imageService := mustNewImageService(t)
	socketService := services.NewSocketService()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("title", "Invalid tag"); err != nil {
		t.Fatalf("WriteField(title) error = %v", err)
	}
	if err := writer.WriteField("tags", "good-tag,bad tag!"); err != nil {
		t.Fatalf("WriteField(tags) error = %v", err)
	}
	part, err := writer.CreateFormFile("file", "pixel.png")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write(mustPNGBytes(t, 10, 10)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/uploads", UploadImage(imageService, socketService))

	req := httptest.NewRequest("POST", "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := performRequest(t, app, req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var errResp models.ErrorResponse
	decodeJSON(t, resp, &errResp)

	if errResp.Message != `tag must use only letters, numbers, and hyphens: "bad tag!"` {
		t.Fatalf("Message = %q, want invalid tag error", errResp.Message)
	}
}

func TestUploadImageRejectsTagLongerThan24Characters(t *testing.T) {
	restore := chdirToBackendRoot(t)
	defer restore()

	imageService := mustNewImageService(t)
	socketService := services.NewSocketService()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("title", "Long tag"); err != nil {
		t.Fatalf("WriteField(title) error = %v", err)
	}
	if err := writer.WriteField("tags", "this-tag-is-way-too-long-for-the-limit"); err != nil {
		t.Fatalf("WriteField(tags) error = %v", err)
	}
	part, err := writer.CreateFormFile("file", "pixel.png")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write(mustPNGBytes(t, 10, 10)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/uploads", UploadImage(imageService, socketService))

	req := httptest.NewRequest("POST", "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := performRequest(t, app, req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var errResp models.ErrorResponse
	decodeJSON(t, resp, &errResp)

	if errResp.Message != `tag must be 24 characters or fewer: "this-tag-is-way-too-long-for-the-limit"` {
		t.Fatalf("Message = %q, want oversized tag error", errResp.Message)
	}
}

func mustNewImageService(t *testing.T) *services.ImageService {
	t.Helper()

	imageService, err := services.NewImageService(config.Config{
		MaxFileSize:       10 * 1024 * 1024,
		MaxImageDimension: 1920,
		AllowedMimeTypes:  []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/avif"},
	})
	if err != nil {
		t.Fatalf("NewImageService() error = %v", err)
	}

	return imageService
}

func chdirToBackendRoot(t *testing.T) func() {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve current test file path")
	}

	return chdir(t, filepath.Dir(filepath.Dir(filename)))
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", dir, err)
	}

	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working directory error = %v", err)
		}
	}
}

func performRequest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	t.Helper()

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
}

func mustPNGBytes(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{R: 220, G: uint8(x % 255), B: uint8(y % 255), A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode() error = %v", err)
	}

	return buf.Bytes()
}
