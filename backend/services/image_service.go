package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
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
)

var (
	ErrTooManyTags = errors.New("maximum 5 tags allowed")
	ErrTagTooLong  = errors.New("tag must be 24 characters or fewer")
	ErrInvalidTag  = errors.New("tag must use only letters, numbers, and hyphens")
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
func (imageService *ImageService) Upload(file *multipart.FileHeader, title string, tags []string, contentType string) (models.Image, error) {
	src, err := file.Open()
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to read uploaded file: %w", err)
	}

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

func (imageService *ImageService) ParseTags(str string) ([]string, error) {
	if str == "" {
		return []string{}, nil
	}
	parts := strings.Split(str, ",")
	// Keep tags unique so repeated values do not bypass the max-tag limit or create duplicate pills in the UI.
	seen := make(map[string]bool)
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := normalizeTag(part)
		if tag == "" {
			continue
		}
		if len([]rune(tag)) > 24 {

			return nil, fmt.Errorf("%w: %q", ErrTagTooLong, tag)
		}
		if !isValidTag(tag) {
			return nil, fmt.Errorf("%w: %q", ErrInvalidTag, tag)
		}
		if seen[tag] {
			continue
		}
		if len(tags) >= 5 {
			return nil, ErrTooManyTags
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	return tags, nil
}

func (imageService *ImageService) ResolveImageType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %w", err)
	}

	contentType, err := detectImageType(data)
	if err != nil {
		return "", err
	}
	if !imageService.IsValidImageType(contentType) {
		return "", fmt.Errorf("unsupported image content type: %s", contentType)
	}

	return contentType, nil
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

func detectImageType(data []byte) (string, error) {
	if _, err := avif.DecodeConfig(bytes.NewReader(data)); err == nil {
		return "image/avif", nil
	}

	if _, err := webp.DecodeConfig(bytes.NewReader(data)); err == nil {
		return "image/webp", nil
	}

	if _, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		switch format {
		case "jpeg":
			return "image/jpeg", nil
		case "png":
			return "image/png", nil
		case "gif":
			return "image/gif", nil
		}
	}

	return "", fmt.Errorf("unsupported or invalid image data")
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

func normalizeTag(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func isValidTag(tag string) bool {
	for _, r := range tag {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			continue
		}
		return false
	}

	return true
}
