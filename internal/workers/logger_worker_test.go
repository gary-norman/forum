package workers

import (
	"context"
	"testing"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

// TestNewLoggerPool tests pool creation
func TestNewLoggerPool(t *testing.T) {
	pool := NewLoggerPool(3, 100, nil)
	if pool == nil {
		t.Fatal("NewLoggerPool returned nil")
	}

	if pool.workers != 3 {
		t.Errorf("expected 3 workers, got %d", pool.workers)
	}

	if cap(pool.logs) != 100 {
		t.Errorf("expected logs channel buffer of 100, got %d", cap(pool.logs))
	}

	if pool.shutdownCh == nil {
		t.Error("shutdown channel was not initialized")
	}
}

// TestLoggerPoolStartAndShutdown tests starting and stopping the pool
func TestLoggerPoolStartAndShutdown(t *testing.T) {
	pool := NewLoggerPool(2, 50, nil)

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

// TestSubmitAndProcessLog tests submitting logs to the pool
func TestSubmitAndProcessLog(t *testing.T) {
	pool := NewLoggerPool(2, 10, nil)
	pool.Start()

	entry := LogEntry{
		Type: "request",
		RequestLog: &models.RequestLog{
			Timestamp:  time.Now(),
			Method:     "GET",
			Path:       "/test",
			StatusCode: 200,
			Duration:   150,
		},
	}

	err := pool.Submit(entry)
	if err != nil {
		t.Errorf("failed to submit log: %v", err)
	}

	// Give worker time to process
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestMultipleLogs tests processing multiple log entries concurrently
func TestMultipleLogs(t *testing.T) {
	pool := NewLoggerPool(3, 50, nil)
	pool.Start()

	// Submit 20 log entries
	for i := 0; i < 20; i++ {
		entry := LogEntry{
			Type: "request",
			RequestLog: &models.RequestLog{
				Timestamp:  time.Now(),
				Method:     "POST",
				Path:       "/api/test",
				StatusCode: 201,
				Duration:   int64(i * 10),
			},
		}

		err := pool.Submit(entry)
		if err != nil {
			t.Errorf("failed to submit log %d: %v", i, err)
		}
	}

	// Give workers time to process all logs
	time.Sleep(500 * time.Millisecond)

	// Verify queue is empty
	stats := pool.Stats()
	if stats.QueueLength != 0 {
		t.Errorf("expected 0 queued logs, got %d", stats.QueueLength)
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestLoggerSubmitAfterShutdown tests that submitting after shutdown fails
func TestLoggerSubmitAfterShutdown(t *testing.T) {
	pool := NewLoggerPool(2, 10, nil)
	pool.Start()

	// Shutdown immediately
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)

	// Try to submit after shutdown
	entry := LogEntry{
		Type: "error",
		ErrorLog: &models.ErrorLog{
			Timestamp: time.Now(),
			Level:     models.LogLevelError,
			Message:   "Test error",
		},
	}

	err := pool.Submit(entry)
	if err == nil {
		t.Error("expected error when submitting after shutdown, got nil")
	}
}

// TestLoggerStats tests pool statistics
func TestLoggerStats(t *testing.T) {
	pool := NewLoggerPool(3, 10, nil)

	stats := pool.Stats()

	if stats.Workers != 3 {
		t.Errorf("expected 3 workers in stats, got %v", stats.Workers)
	}

	if stats.QueueLength != 0 {
		t.Errorf("expected 0 queued logs initially, got %v", stats.QueueLength)
	}

	if stats.QueueCapacity != 10 {
		t.Errorf("expected queue capacity of 10, got %v", stats.QueueCapacity)
	}

	// Submit some logs
	pool.Start()
	for i := 0; i < 5; i++ {
		entry := LogEntry{
			Type: "metric",
			SystemMetric: &models.SystemMetric{
				Timestamp:   time.Now(),
				MetricType:  models.MetricTypeMemory,
				MetricName:  "heap_usage",
				MetricValue: 1024.5,
				Unit:        "MB",
			},
		}
		pool.Submit(entry)
	}

	// Check stats immediately (some logs should be queued or processing)
	stats = pool.Stats()
	if stats.DroppedLogs < 0 {
		t.Errorf("unexpected dropped logs count: %d", stats.DroppedLogs)
	}

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}

// TestHelperMethods tests the convenience methods
func TestHelperMethods(t *testing.T) {
	pool := NewLoggerPool(2, 10, nil)
	pool.Start()

	// Test LogRequest
	err := pool.LogRequest(models.RequestLog{
		Timestamp:  time.Now(),
		Method:     "GET",
		Path:       "/",
		StatusCode: 200,
		Duration:   50,
	})
	if err != nil {
		t.Errorf("LogRequest failed: %v", err)
	}

	// Test LogError
	err = pool.LogError(models.ErrorLog{
		Timestamp: time.Now(),
		Level:     models.LogLevelWarn,
		Message:   "Test warning",
	})
	if err != nil {
		t.Errorf("LogError failed: %v", err)
	}

	// Test LogMetric
	err = pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeHealthCheck,
		MetricName:  "db_ping",
		MetricValue: 5.2,
		Unit:        "ms",
	})
	if err != nil {
		t.Errorf("LogMetric failed: %v", err)
	}

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}
