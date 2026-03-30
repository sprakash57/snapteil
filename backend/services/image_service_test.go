package services

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sprakash57/snapteil/backend/config"
	"github.com/sprakash57/snapteil/backend/models"
)

func TestParseTagsSanitizesAndDeduplicates(t *testing.T) {
	imageService := &ImageService{}

	got := imageService.ParseTags(" Nature , city!, city, street-fight, ----, this-tag-is-way-too-long-for-the-limit ")
	want := []string{"nature", "city", "street-fight"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseTags() = %v, want %v", got, want)
	}
}

func TestIsValidImageTypeNormalizesMimeParameters(t *testing.T) {
	imageService := &ImageService{
		cfg: config.Config{
			AllowedMimeTypes: []string{"image/png", "image/jpeg"},
		},
	}

	if !imageService.IsValidImageType("image/png; charset=binary") {
		t.Fatal("expected image/png with parameters to be accepted")
	}

	if imageService.IsValidImageType("image/svg+xml") {
		t.Fatal("expected image/svg+xml to be rejected")
	}
}

func TestAddMaintainsNewestFirstOrder(t *testing.T) {
	base := time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC)
	imageService := &ImageService{
		images: []models.Image{
			{ID: "2", CreatedAt: base.Add(-1 * time.Minute)},
			{ID: "3", CreatedAt: base.Add(-3 * time.Minute)},
		},
	}

	imageService.Add(models.Image{ID: "1", CreatedAt: base})
	imageService.Add(models.Image{ID: "4", CreatedAt: base.Add(-2 * time.Minute)})

	got := []string{
		imageService.images[0].ID,
		imageService.images[1].ID,
		imageService.images[2].ID,
		imageService.images[3].ID,
	}
	want := []string{"1", "2", "4", "3"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("image order = %v, want %v", got, want)
	}
}

func TestGetPaginatedFiltersByTags(t *testing.T) {
	base := time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC)
	imageService := &ImageService{
		images: []models.Image{
			{ID: "1", Tags: []string{"got"}, CreatedAt: base},
			{ID: "2", Tags: []string{"naruto"}, CreatedAt: base.Add(-1 * time.Minute)},
			{ID: "3", Tags: []string{"skyrim"}, CreatedAt: base.Add(-2 * time.Minute)},
		},
	}

	resp := imageService.GetPaginated(1, 1, []string{"naruto", "skyrim"})

	if resp.Total != 2 {
		t.Fatalf("Total = %d, want 2", resp.Total)
	}
	if !resp.HasMore {
		t.Fatal("expected HasMore to be true for first filtered page")
	}
	if len(resp.Images) != 1 || resp.Images[0].ID != "2" {
		t.Fatalf("Images = %#v, want first filtered image ID 2", resp.Images)
	}
}

func TestResolveImageTypeDetectsPngAndRejectsInvalidData(t *testing.T) {
	imageService := &ImageService{
		cfg: config.Config{
			AllowedMimeTypes: []string{"image/png", "image/jpeg"},
		},
	}

	t.Run("png", func(t *testing.T) {
		fileHeader := makeMultipartFileHeader(t, "pixel.png", mustPNGBytes(t, 10, 10))

		got, err := imageService.ResolveImageType(fileHeader)
		if err != nil {
			t.Fatalf("ResolveImageType() error = %v", err)
		}
		if got != "image/png" {
			t.Fatalf("ResolveImageType() = %q, want image/png", got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		fileHeader := makeMultipartFileHeader(t, "notes.txt", []byte("not an image"))

		if _, err := imageService.ResolveImageType(fileHeader); err == nil {
			t.Fatal("expected invalid file to be rejected")
		}
	})
}

func TestUploadPreservesOriginalFormatForSmallPng(t *testing.T) {
	restore := chdirForTest(t, t.TempDir())
	defer restore()

	imageService := &ImageService{
		cfg: config.Config{
			MaxFileSize:       1024 * 1024,
			MaxImageDimension: 500,
			AllowedMimeTypes:  []string{"image/png"},
		},
	}

	original := mustPNGBytes(t, 20, 10)
	fileHeader := makeMultipartFileHeader(t, "small.png", original)

	record, err := imageService.Upload(fileHeader, "Small", []string{"tag"}, "image/png")
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if !strings.HasSuffix(record.Filename, ".png") {
		t.Fatalf("Filename = %q, want .png suffix", record.Filename)
	}
	if record.Width != 20 || record.Height != 10 {
		t.Fatalf("dimensions = %dx%d, want 20x10", record.Width, record.Height)
	}

	storedPath := filepath.Join("uploads", record.Filename)
	stored, err := os.ReadFile(storedPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", storedPath, err)
	}

	if !bytes.Equal(stored, original) {
		t.Fatal("expected upload to keep original PNG bytes when no resize is needed")
	}
}

func TestUploadResizesLargePngAndKeepsFormat(t *testing.T) {
	restore := chdirForTest(t, t.TempDir())
	defer restore()

	imageService := &ImageService{
		cfg: config.Config{
			MaxFileSize:       2 * 1024 * 1024,
			MaxImageDimension: 100,
			AllowedMimeTypes:  []string{"image/png"},
		},
	}

	fileHeader := makeMultipartFileHeader(t, "large.png", mustPNGBytes(t, 200, 100))

	record, err := imageService.Upload(fileHeader, "Large", []string{"tag"}, "image/png")
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if !strings.HasSuffix(record.Filename, ".png") {
		t.Fatalf("Filename = %q, want .png suffix", record.Filename)
	}
	if record.Width != 100 || record.Height != 50 {
		t.Fatalf("dimensions = %dx%d, want 100x50", record.Width, record.Height)
	}

	storedPath := filepath.Join("uploads", record.Filename)
	stored, err := os.ReadFile(storedPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", storedPath, err)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(stored))
	if err != nil {
		t.Fatalf("DecodeConfig() error = %v", err)
	}
	if format != "png" {
		t.Fatalf("stored format = %q, want png", format)
	}
	if cfg.Width != 100 || cfg.Height != 50 {
		t.Fatalf("stored dimensions = %dx%d, want 100x50", cfg.Width, cfg.Height)
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

func makeMultipartFileHeader(t *testing.T, filename string, data []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(body.Len())); err != nil {
		t.Fatalf("ParseMultipartForm() error = %v", err)
	}

	fileHeaders := req.MultipartForm.File["file"]
	if len(fileHeaders) != 1 {
		t.Fatalf("expected one multipart file, got %d", len(fileHeaders))
	}

	return fileHeaders[0]
}

func chdirForTest(t *testing.T, dir string) func() {
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
