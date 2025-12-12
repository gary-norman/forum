package workers

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

// TestIntegrationRealImage tests the worker pool with an actual image file
// This test creates a real PNG image, processes it through the worker pool,
// and verifies the processed image is saved correctly
func TestIntegrationRealImage(t *testing.T) {
	// Skip if in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary test image
	tmpFile := filepath.Join(os.TempDir(), "integration-test-image.png")
	createTestImage(t, tmpFile, 200, 150)
	defer os.Remove(tmpFile)

	// Ensure output directory exists
	outputDir := "db/userdata/images/post-images"
	os.MkdirAll(outputDir, 0755)

	// Create and start worker pool (no database for integration test)
	pool := NewImageWorkerPool(2, 10, nil)
	pool.Start()

	// Submit job
	job := ImageJob{
		ID:       "integration-test-job",
		FilePath: tmpFile,
		UserID:   models.NewUUIDField(),
		PostID:   999,
	}

	err := pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Give worker time to process
	time.Sleep(500 * time.Millisecond)

	// Verify the processed image was created
	processedPath := filepath.Join(outputDir, "integration-test-job.png")
	if _, err := os.Stat(processedPath); os.IsNotExist(err) {
		t.Errorf("Processed image not found at %s", processedPath)
	} else {
		// Clean up processed image
		os.Remove(processedPath)
	}

	// Verify temp file was deleted
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Errorf("Temporary file should have been deleted, but still exists at %s", tmpFile)
	}

	// Shutdown pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

// createTestImage creates a simple colored rectangle PNG for testing
func createTestImage(t *testing.T, path string, width, height int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}
