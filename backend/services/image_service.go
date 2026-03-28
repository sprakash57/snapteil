package services

import (
	"encoding/json"
	"os"
	"slices"
	"sort"
	"sync"

	"github.com/sprakash57/snaptiel/backend/models"
)

type ImageService struct {
	mu     sync.RWMutex
	images []models.Image
}

func NewImageService(seedPath string) (*ImageService, error) {
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return nil, err
	}

	var images []models.Image
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}

	return &ImageService{images: images}, nil
}

func (imageService *ImageService) GetPaginated(page, perPage int, tag string) models.PaginatedResponse {
	imageService.mu.RLock()
	defer imageService.mu.RUnlock()

	// Filter by tag if provided, otherwise use all images
	filtered := imageService.images
	if tag != "" {
		filtered = make([]models.Image, 0)
		for _, img := range imageService.images {
			if slices.Contains(img.Tags, tag) {
				filtered = append(filtered, img)
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
