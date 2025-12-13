package workers

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"sync/atomic"

	"github.com/gary-norman/forum/internal/colors"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/sqlite"
)

var (
	loggerColors, _ = colors.UseFlavor("Mocha")
)

// LogEntry represents a log entry to be written to the database
type LogEntry struct {
	Type         string // "request", "error", "metric"
	RequestLog   *models.RequestLog
	ErrorLog     *models.ErrorLog
	SystemMetric *models.SystemMetric
}

// LoggerPool manages a pool of worker goroutines that write logs to the database
type LoggerPool struct {
	logs         chan LogEntry        // Channel for queuing log entries
	workers      int                  // Number of worker goroutines
	wg           sync.WaitGroup       // Wait for workers to finish during shutdown
	shutdownCh   chan struct{}        // Signal to shutdown workers
	isShutdown   atomic.Bool          // Track whether pool has been shut down
	loggingModel *sqlite.LoggingModel // Database model for storing logs
	droppedLogs  atomic.Int64         // Count of dropped logs when queue is full
}

// NewLoggerPool creates a new logger worker pool
// workers: number of concurrent worker goroutines
// queueSize: maximum number of log entries that can be queued
// db: database connection for storing logs (can be nil for tests)
func NewLoggerPool(workers, queueSize int, db *sql.DB) *LoggerPool {
	pool := &LoggerPool{
		logs:       make(chan LogEntry, queueSize),
		workers:    workers,
		shutdownCh: make(chan struct{}),
	}

	// Only create loggingModel if we have a database connection
	if db != nil {
		pool.loggingModel = &sqlite.LoggingModel{DB: db}
	}

	return pool
}

// Start starts the logger worker pool
func (pool *LoggerPool) Start() {
	for i := 0; i < pool.workers; i++ {
		pool.wg.Add(1)
		go func(workerID int) {
			defer pool.wg.Done()
			for {
				select {
				case entry := <-pool.logs:
					pool.writeLog(entry, workerID)
				case <-pool.shutdownCh:
					return
				}
			}
		}(i)
	}
	log.Printf(loggerColors.Blue + "[LoggerPool] Started with %d workers" + loggerColors.Reset + "\n", pool.workers)
}

// TODO(human): Exercise 4 Part 1 - Implement Submit method
//
// Instructions:
// This is very similar to the ImageWorkerPool.Submit() you implemented!
//
// 1. Check if pool is shut down using pool.isShutdown.Load()
//    - If shut down, return error "logger pool is shut down"
//
// 2. Use select with default to submit entry non-blocking:
//    select {
//    case pool.logs <- entry:
//        return nil
//    default:
//        // Queue is full - increment dropped counter and return error
//        pool.droppedLogs.Add(1)
//        return fmt.Errorf("logger pool queue is full (dropped: %d)", pool.droppedLogs.Load())
//    }
//
// Why non-blocking? Logging should NEVER slow down your application!
// If the queue is full, it's better to drop the log than block the HTTP handler.

// Submit submits a log entry to the worker pool
// Returns an error if the queue is full or pool is shut down
func (pool *LoggerPool) Submit(entry LogEntry) error {
	// TODO(human): Implement non-blocking submission with shutdown check
	return nil
}

// TODO(human): Exercise 4 Part 2 - Implement Shutdown method
//
// Instructions:
// This is identical to ImageWorkerPool.Shutdown() pattern!
//
// 1. Set shutdown flag: pool.isShutdown.Store(true)
// 2. Close shutdown channel: close(pool.shutdownCh)
// 3. Create done channel: done := make(chan struct{})
// 4. Start goroutine:
//    go func() {
//        pool.wg.Wait()  // Wait for all workers
//        close(done)     // Signal completion
//    }()
// 5. Select between done and context timeout:
//    select {
//    case <-done:
//        return nil
//    case <-ctx.Done():
//        return ctx.Err()
//    }

// Shutdown gracefully shuts down the logger pool
// Waits for all workers to finish writing pending logs
// Returns error if context times out before workers finish
func (pool *LoggerPool) Shutdown(ctx context.Context) error {
	// TODO(human): Implement graceful shutdown with context timeout
	return nil
}

// writeLog processes a single log entry and writes it to the database
func (pool *LoggerPool) writeLog(entry LogEntry, workerID int) {
	if pool.loggingModel == nil {
		return // No database in test environment
	}

	var err error

	switch entry.Type {
	case "request":
		if entry.RequestLog != nil {
			err = pool.loggingModel.InsertRequestLog(*entry.RequestLog)
		}
	case "error":
		if entry.ErrorLog != nil {
			err = pool.loggingModel.InsertErrorLog(*entry.ErrorLog)
		}
	case "metric":
		if entry.SystemMetric != nil {
			err = pool.loggingModel.InsertSystemMetric(*entry.SystemMetric)
		}
	}

	if err != nil {
		// Log to console if database write fails (don't want to create infinite loop!)
		log.Printf(loggerColors.Red+"[LoggerWorker %d] Failed to write %s log: %v"+loggerColors.Reset+"\n",
			workerID, entry.Type, err)
	}
}

// Stats returns statistics about the logger pool
type LoggerStats struct {
	Workers      int
	QueueLength  int
	QueueCapacity int
	QueueUsage   float64 // Percentage of queue used (0.0 to 1.0)
	DroppedLogs  int64
}

func (pool *LoggerPool) Stats() LoggerStats {
	queueLen := len(pool.logs)
	queueCap := cap(pool.logs)
	usage := 0.0
	if queueCap > 0 {
		usage = float64(queueLen) / float64(queueCap)
	}

	return LoggerStats{
		Workers:       pool.workers,
		QueueLength:   queueLen,
		QueueCapacity: queueCap,
		QueueUsage:    usage,
		DroppedLogs:   pool.droppedLogs.Load(),
	}
}

// Helper methods for easy log submission

// LogRequest submits a request log entry
func (pool *LoggerPool) LogRequest(log models.RequestLog) error {
	return pool.Submit(LogEntry{
		Type:       "request",
		RequestLog: &log,
	})
}

// LogError submits an error log entry
func (pool *LoggerPool) LogError(log models.ErrorLog) error {
	return pool.Submit(LogEntry{
		Type:     "error",
		ErrorLog: &log,
	})
}

// LogMetric submits a system metric entry
func (pool *LoggerPool) LogMetric(metric models.SystemMetric) error {
	return pool.Submit(LogEntry{
		Type:         "metric",
		SystemMetric: &metric,
	})
}
