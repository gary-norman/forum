package workers

import (
	"context"
	"database/sql"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gary-norman/forum/internal/models"
)

// TestDatabaseIntegration tests the complete flow with a real database
func TestDatabaseIntegration(t *testing.T) {
	// Skip if in short mode
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	// Create in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create Images table with Path column
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Images (
		ID INTEGER PRIMARY KEY,
		Created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		Updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		AuthorID BLOB NOT NULL,
		PostID INTEGER NOT NULL,
		Path TEXT NOT NULL DEFAULT ''
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatalf("Failed to create Images table: %v", err)
	}

	// Create a temporary test image
	tmpFile := filepath.Join(os.TempDir(), "db-integration-test-image.png")
	createTestImageForDB(t, tmpFile, 100, 100)
	defer os.Remove(tmpFile)

	// Ensure output directory exists
	outputDir := "db/userdata/images/post-images"
	os.MkdirAll(outputDir, 0755)

	// Create and start worker pool WITH database
	pool := NewImageWorkerPool(1, 5, db)
	pool.Start()

	// Submit job
	testUserID := models.NewUUIDField()
	testPostID := int64(42)

	job := ImageJob{
		ID:       "db-integration-test",
		FilePath: tmpFile,
		UserID:   testUserID,
		PostID:   testPostID,
	}

	err = pool.Submit(job)
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Give worker time to process
	time.Sleep(500 * time.Millisecond)

	// Verify the processed image was created
	processedPath := filepath.Join(outputDir, "db-integration-test.png")
	if _, err := os.Stat(processedPath); os.IsNotExist(err) {
		t.Errorf("Processed image not found at %s", processedPath)
	} else {
		// Clean up processed image
		defer os.Remove(processedPath)
	}

	// Verify database record was created
	var imageID int64
	var dbPath string
	var dbAuthorID []byte
	var dbPostID int64

	query := "SELECT ID, Path, AuthorID, PostID FROM Images WHERE Path = ?"
	err = db.QueryRow(query, processedPath).Scan(&imageID, &dbPath, &dbAuthorID, &dbPostID)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	// Verify data
	if dbPath != processedPath {
		t.Errorf("Expected path %s, got %s", processedPath, dbPath)
	}

	if dbPostID != testPostID {
		t.Errorf("Expected PostID %d, got %d", testPostID, dbPostID)
	}

	// Verify AuthorID matches (compare bytes)
	expectedAuthorIDBytes, _ := testUserID.Value()
	expectedBytes := expectedAuthorIDBytes.([]byte)

	if len(dbAuthorID) != len(expectedBytes) {
		t.Errorf("AuthorID length mismatch: expected %d bytes, got %d bytes", len(expectedBytes), len(dbAuthorID))
	} else {
		for i := range expectedBytes {
			if dbAuthorID[i] != expectedBytes[i] {
				t.Errorf("AuthorID mismatch at byte %d: expected %x, got %x", i, expectedBytes[i], dbAuthorID[i])
				break
			}
		}
	}

	t.Logf("✓ Database integration successful! Image ID: %d, Path: %s", imageID, dbPath)

	// Shutdown pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

// createTestImageForDB creates a simple colored rectangle PNG for database testing
func createTestImageForDB(t *testing.T, path string, width, height int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a solid green square
	green := color.RGBA{0, 255, 0, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, green)
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
