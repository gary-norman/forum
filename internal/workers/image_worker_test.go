package workers

import (
	"context"
	"testing"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

// TestNewImageWorkerPool tests pool creation
func TestNewImageWorkerPool(t *testing.T) {
	pool := NewImageWorkerPool(3, 10)
	if pool == nil {
		t.Fatal("NewImageWorkerPool returned nil")
	}

	if pool.workers != 3 {
		t.Errorf("expected 3 workers, got %d", pool.workers)
	}

	if cap(pool.jobs) != 10 {
		t.Errorf("expected jobs channel buffer of 10, got %d", cap(pool.jobs))
	}

	if pool.shutdownCh == nil {
		t.Error("shutdown channel was not initialized")
	}
}

// TestWorkerPoolStartAndShutdown tests starting and stopping the pool
func TestWorkerPoolStartAndShutdown(t *testing.T) {
	pool := NewImageWorkerPool(2, 5)

	// Start the pool
	pool.Start()

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.Shutdown(ctx)
	if err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestSubmitAndProcessJob tests submitting jobs to the pool
func TestSubmitAndProcessJob(t *testing.T) {
	pool := NewImageWorkerPool(2, 10)
	pool.Start()

	job := ImageJob{
		ID:       "test-job-1",
		FilePath: "/tmp/test.jpg",
		UserID:   models.NewUUIDField(),
		PostID:   1,
	}

	err := pool.Submit(job)
	if err != nil {
		t.Errorf("failed to submit job: %v", err)
	}

	// Give worker time to process
	time.Sleep(1 * time.Second)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestMultipleJobs tests processing multiple jobs concurrently
func TestMultipleJobs(t *testing.T) {
	pool := NewImageWorkerPool(3, 20)
	pool.Start()

	// Submit 10 jobs
	for i := 0; i < 10; i++ {
		job := ImageJob{
			ID:       "job-" + string(rune('0'+i)),
			FilePath: "/tmp/image" + string(rune('0'+i)) + ".jpg",
			UserID:   models.NewUUIDField(),
			PostID:   int64(i + 1),
		}

		err := pool.Submit(job)
		if err != nil {
			t.Errorf("failed to submit job %d: %v", i, err)
		}
	}

	// Give workers time to process all jobs
	// With 3 workers and 500ms per job, this should take ~2 seconds
	time.Sleep(3 * time.Second)

	// Verify queue is empty
	stats := pool.Stats()
	queuedJobs := stats.QueueLength
	if queuedJobs != 0 {
		t.Errorf("expected 0 queued jobs, got %d", queuedJobs)
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestSubmitAfterShutdown tests that submitting after shutdown fails
func TestSubmitAfterShutdown(t *testing.T) {
	pool := NewImageWorkerPool(2, 5)
	pool.Start()

	// Shutdown immediately
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)

	// Try to submit after shutdown
	job := ImageJob{
		ID:       "test-job-after-shutdown",
		FilePath: "/tmp/test.jpg",
		UserID:   models.NewUUIDField(),
		PostID:   1,
	}

	err := pool.Submit(job)
	if err == nil {
		t.Error("expected error when submitting after shutdown, got nil")
	}
}

// TestBufferedChannels tests that buffered channels prevent blocking
func TestBufferedChannels(t *testing.T) {
	// Create pool with buffer size of 5
	pool := NewImageWorkerPool(1, 5)

	// Don't start workers yet - jobs will queue in buffer

	// Submit 5 jobs (should not block due to buffer)
	for i := 0; i < 5; i++ {
		job := ImageJob{
			ID:       "buffered-job-" + string(rune('0'+i)),
			FilePath: "/tmp/image.jpg",
			UserID:   models.NewUUIDField(),
			PostID:   1,
		}

		// This should complete immediately without blocking
		done := make(chan error, 1)
		go func() {
			done <- pool.Submit(job)
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Errorf("job %d failed to submit: %v", i, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("job %d submission blocked (buffer should prevent this)", i)
		}
	}

	// Verify all 5 jobs are queued
	stats := pool.Stats()
	queuedJobs := stats.QueueLength
	if queuedJobs != 5 {
		t.Errorf("expected 5 queued jobs, got %d", queuedJobs)
	}

	// Now start workers and let them process
	pool.Start()
	time.Sleep(3 * time.Second)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestConcurrentProcessing verifies parallel execution
func TestConcurrentProcessing(t *testing.T) {
	// Create pool with 3 workers
	pool := NewImageWorkerPool(3, 20)
	pool.Start()

	// Submit 6 jobs
	for i := 0; i < 6; i++ {
		job := ImageJob{
			ID:       "concurrent-job-" + string(rune('0'+i)),
			FilePath: "/tmp/image.jpg",
			UserID:   models.NewUUIDField(),
			PostID:   1,
		}
		err := pool.Submit(job)
		if err != nil {
			t.Errorf("failed to submit job %d: %v", i, err)
		}
	}

	// Give workers time to process (jobs will fail validation but that's ok)
	time.Sleep(100 * time.Millisecond)

	// Verify queue is empty (all jobs picked up by workers)
	stats := pool.Stats()
	if stats.QueueLength != 0 {
		t.Errorf("expected empty queue after processing, got %d queued", stats.QueueLength)
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestShutdownTimeout tests shutdown with a timeout
func TestShutdownTimeout(t *testing.T) {
	pool := NewImageWorkerPool(2, 10)
	pool.Start()

	// Submit some jobs (they'll process quickly since validation fails)
	for i := 0; i < 4; i++ {
		job := ImageJob{
			ID:       "quick-job-" + string(rune('0'+i)),
			FilePath: "/tmp/image.jpg",
			UserID:   models.NewUUIDField(),
			PostID:   1,
		}
		pool.Submit(job)
	}

	// Give workers time to pick up jobs
	time.Sleep(50 * time.Millisecond)

	// Shutdown should succeed even with short timeout since jobs complete quickly
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := pool.Shutdown(ctx)
	if err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestGetStats tests pool statistics
func TestGetStats(t *testing.T) {
	pool := NewImageWorkerPool(3, 10)

	stats := pool.Stats()

	if stats.Workers != 3 {
		t.Errorf("expected 3 workers in stats, got %v", stats.Workers)
	}

	if stats.QueueLength != 0 {
		t.Errorf("expected 0 queued jobs initially, got %v", stats.QueueLength)
	}

	// Submit some jobs
	pool.Start()
	for i := 0; i < 5; i++ {
		job := ImageJob{
			ID:       "stats-job-" + string(rune('0'+i)),
			FilePath: "/tmp/image.jpg",
			UserID:   models.NewUUIDField(),
			PostID:   1,
		}
		pool.Submit(job)
	}

	// Check stats immediately (some jobs should be queued)
	stats = pool.Stats()
	queuedJobs := stats.QueueLength
	if queuedJobs < 0 || queuedJobs > 5 {
		t.Errorf("unexpected queued jobs count: %d", queuedJobs)
	}

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}
