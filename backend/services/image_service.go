package services

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/sprakash57/snapteil/backend/config"
	"github.com/sprakash57/snapteil/backend/models"
	_ "golang.org/x/image/webp"
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

// Upload decodes, normalizes, and returns the saved uploaded image file.
func (imageService *ImageService) Upload(file *multipart.FileHeader, title string, tags []string) (models.Image, error) {
	src, err := file.Open()
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to decode image: %w", err)
	}

	img = imageService.normalizeImage(img)
	bounds := img.Bounds()

	id := uuid.New().String()
	filename := id + ".jpg"
	filePath := filepath.Join("./uploads", filename)

	if err := os.MkdirAll("./uploads", 0o755); err != nil {
		return models.Image{}, fmt.Errorf("failed to ensure upload directory: %w", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to save image: %w", err)
	}
	defer outFile.Close()
	// Encode as JPEG with quality 85 to balance size and quality
	if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 85}); err != nil {
		return models.Image{}, fmt.Errorf("failed to encode image: %w", err)
	}

	record := models.Image{
		ID:        id,
		Title:     title,
		Tags:      tags,
		Filename:  filename,
		URL:       "/uploads/" + filename,
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
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
	return slices.Contains(imageService.cfg.AllowedMimeTypes, contentType)
}

func (imageService *ImageService) MaxFileSize() int64 {
	return imageService.cfg.MaxFileSize
}

func (imageService *ImageService) normalizeImage(img image.Image) image.Image {
	maxDim := imageService.cfg.MaxImageDimension
	return imaging.Fit(img, maxDim, maxDim, imaging.Lanczos)
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
