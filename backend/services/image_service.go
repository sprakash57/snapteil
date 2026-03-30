package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/google/uuid"
	"github.com/sprakash57/snapteil/backend/config"
	"github.com/sprakash57/snapteil/backend/models"
	"github.com/srwiley/oksvg"
)

type ImageService struct {
	mu     sync.RWMutex
	images []models.Image
	cfg    config.Config
}

func NewImageService(cfg config.Config) (*ImageService, error) {
	data, err := os.ReadFile("./data/seed.json")
	if err != nil {
		return nil, err
	}

	var images []models.Image
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].CreatedAt.After(images[j].CreatedAt)
	})

	log.Printf("Loaded %d seed images from ./data/seed.json", len(images))

	return &ImageService{images: images, cfg: cfg}, nil
}

func (imageService *ImageService) GetPaginated(page, perPage int, tags []string) models.PaginatedResponse {
	imageService.mu.RLock()
	defer imageService.mu.RUnlock()

	if len(tags) > 0 {
		filtered := make([]models.Image, 0)
		for _, img := range imageService.images {
			if matchesAnyTag(img.Tags, tags) {
				filtered = append(filtered, img)
			}
		}

		return paginatedResponse(filtered, page, perPage)
	}

	return paginatedResponse(imageService.images, page, perPage)
}

func (imageService *ImageService) Add(img models.Image) {
	imageService.mu.Lock()
	defer imageService.mu.Unlock()

	insertAt := sort.Search(len(imageService.images), func(i int) bool {
		return !imageService.images[i].CreatedAt.After(img.CreatedAt)
	})

	imageService.images = append(imageService.images, models.Image{})
	copy(imageService.images[insertAt+1:], imageService.images[insertAt:])
	imageService.images[insertAt] = img
}

// Upload stores the image in its original format and only re-encodes if resizing is required.
func (imageService *ImageService) Upload(file *multipart.FileHeader, title string, tags []string) (models.Image, error) {
	src, err := file.Open()
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to read uploaded file: %w", err)
	}

	contentType := normalizeContentType(file.Header.Get("Content-Type"))

	encoded, width, height, err := imageService.normalizeUpload(data, contentType)
	if err != nil {
		return models.Image{}, err
	}

	id := uuid.New().String()
	filename := id + extensionForContentType(contentType)
	filePath := filepath.Join("./uploads", filename)

	if err := os.MkdirAll("./uploads", 0o755); err != nil {
		return models.Image{}, fmt.Errorf("failed to ensure upload directory: %w", err)
	}

	if err := os.WriteFile(filePath, encoded, 0o644); err != nil {
		return models.Image{}, fmt.Errorf("failed to save image: %w", err)
	}

	record := models.Image{
		ID:        id,
		Title:     title,
		Tags:      tags,
		Filename:  filename,
		URL:       "/uploads/" + filename,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now(),
	}

	imageService.Add(record)
	return record, nil
}

func (imageService *ImageService) ParseTags(str string) []string {
	if str == "" {
		return []string{}
	}
	parts := strings.Split(str, ",")
	seen := make(map[string]bool)
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := sanitizeTag(part)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	return tags
}

func (imageService *ImageService) IsValidImageType(contentType string) bool {
	return slices.Contains(imageService.cfg.AllowedMimeTypes, normalizeContentType(contentType))
}

func (imageService *ImageService) MaxFileSize() int64 {
	return imageService.cfg.MaxFileSize
}

func (imageService *ImageService) normalizeImage(img image.Image) image.Image {
	maxDim := imageService.cfg.MaxImageDimension
	return imaging.Fit(img, maxDim, maxDim, imaging.Lanczos)
}

func (imageService *ImageService) normalizeUpload(data []byte, contentType string) ([]byte, int, int, error) {
	if contentType == "image/svg+xml" {
		width, height, err := imageService.decodeSVGDimensions(data)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to decode image: %w", err)
		}

		width, height = imageService.normalizeDimensions(width, height)
		return data, width, height, nil
	}

	img, err := decodeImage(bytes.NewReader(data), contentType)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	originalBounds := img.Bounds()
	normalized := imageService.normalizeImage(img)
	bounds := normalized.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == originalBounds.Dx() && height == originalBounds.Dy() {
		return data, width, height, nil
	}

	var encoded bytes.Buffer
	if err := encodeImage(&encoded, normalized, contentType); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to encode image: %w", err)
	}

	return encoded.Bytes(), width, height, nil
}

func (imageService *ImageService) decodeSVGDimensions(data []byte) (int, int, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		return 0, 0, err
	}

	width := int(math.Round(icon.ViewBox.W))
	height := int(math.Round(icon.ViewBox.H))
	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("svg has invalid dimensions")
	}

	return width, height, nil
}

func (imageService *ImageService) normalizeDimensions(width, height int) (int, int) {
	maxDim := imageService.cfg.MaxImageDimension
	if width <= maxDim && height <= maxDim {
		return width, height
	}

	scale := math.Min(float64(maxDim)/float64(width), float64(maxDim)/float64(height))
	width = int(math.Round(float64(width) * scale))
	height = int(math.Round(float64(height) * scale))

	return max(width, 1), max(height, 1)
}

func decodeImage(r io.Reader, contentType string) (image.Image, error) {
	switch contentType {
	case "image/avif":
		return avif.Decode(r)
	case "image/webp":
		return webp.Decode(r)
	default:
		img, _, err := image.Decode(r)
		return img, err
	}
}

func encodeImage(w io.Writer, img image.Image, contentType string) error {
	switch contentType {
	case "image/jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	case "image/png":
		return png.Encode(w, img)
	case "image/gif":
		return gif.Encode(w, img, nil)
	case "image/webp":
		return webp.Encode(w, img)
	case "image/avif":
		return avif.Encode(w, img)
	default:
		return fmt.Errorf("unsupported image content type: %s", contentType)
	}
}

func extensionForContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/avif":
		return ".avif"
	case "image/svg+xml":
		return ".svg"
	default:
		if extensions, err := mime.ExtensionsByType(contentType); err == nil && len(extensions) > 0 {
			return extensions[0]
		}
		return ".img"
	}
}

func normalizeContentType(contentType string) string {
	if contentType == "" {
		return ""
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return strings.TrimSpace(contentType)
	}

	return mediaType
}

func paginatedResponse(images []models.Image, page, perPage int) models.PaginatedResponse {
	total := len(images)
	start := min((page-1)*perPage, total)
	end := min(start+perPage, total)
	pageImages := append([]models.Image(nil), images[start:end]...)

	return models.PaginatedResponse{
		Images:  pageImages,
		Total:   total,
		Page:    page,
		PerPage: perPage,
		HasMore: end < total,
	}
}

func matchesAnyTag(imageTags, filterTags []string) bool {
	for _, tag := range filterTags {
		if slices.Contains(imageTags, tag) {
			return true
		}
	}
	return false
}

func sanitizeTag(raw string) string {
	tag := strings.TrimSpace(strings.ToLower(raw))
	if tag == "" {
		return ""
	}

	var cleaned strings.Builder
	for _, r := range tag {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			cleaned.WriteRune(r)
		}
	}

	tag = strings.Trim(cleaned.String(), "-")
	if tag == "" || len(tag) > 24 {
		return ""
	}

	return tag
}
