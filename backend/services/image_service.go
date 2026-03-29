package services

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
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

	return &ImageService{images: images, cfg: cfg}, nil
}

func (imageService *ImageService) normalizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	maxDim := imageService.cfg.MaxImageDimension

	if width <= maxDim && height <= maxDim {
		return img
	}

	if width > height {
		return imaging.Resize(img, maxDim, 0, imaging.Lanczos)
	}
	return imaging.Resize(img, 0, maxDim, imaging.Lanczos)
}

func (imageService *ImageService) GetPaginated(page, perPage int, tags []string) models.PaginatedResponse {
	imageService.mu.RLock()
	defer imageService.mu.RUnlock()

	// Filter by tags if provided, otherwise use all images
	filtered := imageService.images
	if len(tags) > 0 {
		filtered = make([]models.Image, 0)
		for _, img := range imageService.images {
			for _, tag := range tags {
				if slices.Contains(img.Tags, tag) {
					filtered = append(filtered, img)
					break
				}
			}
		}
	}

	// Sort a copy so we don't mutate the original slice
	sorted := make([]models.Image, len(filtered))
	copy(sorted, filtered)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})

	// Calculate pagination limites
	total := len(sorted)
	start := min((page-1)*perPage, total)
	end := min(start+perPage, total) // ensure end should not go outside of limits

	return models.PaginatedResponse{
		Images:  sorted[start:end], // slice the page window
		Total:   total,
		Page:    page,
		PerPage: perPage,
		HasMore: end < total,
	}
}

func (imageService *ImageService) Add(img models.Image) {
	imageService.mu.Lock()
	defer imageService.mu.Unlock()
	imageService.images = append(imageService.images, img)
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
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(strings.ToLower(part))
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (imageService *ImageService) IsValidImageType(contentType string) bool {
	return slices.Contains(imageService.cfg.AllowedMimeTypes, contentType)
}

func (imageService *ImageService) MaxFileSize() int64 {
	return imageService.cfg.MaxFileSize
}
