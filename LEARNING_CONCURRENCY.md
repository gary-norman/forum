# Go Concurrency: Context and Channels for Server Efficiency

## Overview

This guide will teach you how to use Go's `context` and `channels` to improve your server's efficiency, reliability, and responsiveness.

---

## Part 1: Understanding Context

### What is Context?

Context is Go's way of managing request-scoped values, cancellation signals, and deadlines across API boundaries and goroutines.

**Key Use Cases:**
1. **Cancellation**: Stop expensive operations when the client disconnects
2. **Timeouts**: Prevent operations from running too long
3. **Request-scoped values**: Pass user info, request IDs, etc. through the call stack
4. **Tracing**: Track requests across microservices

### Context in Database Operations

**Problem**: When a user navigates away from a page, their HTTP request is cancelled, but your database query keeps running, wasting resources.

**Solution**: Use context-aware database methods that automatically cancel queries when the context is done.

```go
// ❌ Bad: Query continues even if user disconnects
rows, err := db.Query("SELECT * FROM Posts WHERE ...")

// ✅ Good: Query stops if request is cancelled
rows, err := db.QueryContext(ctx, "SELECT * FROM Posts WHERE ...")
```

### Context Methods

```go
// Create contexts
ctx := context.Background()                              // Root context
ctx := context.TODO()                                    // When you're not sure which context to use
ctx, cancel := context.WithCancel(parentCtx)            // Cancellable
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)  // Auto-cancel after timeout
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))

// Always defer cancel() to prevent context leaks
defer cancel()

// Check if context is cancelled
select {
case <-ctx.Done():
    return ctx.Err() // Returns context.Canceled or context.DeadlineExceeded
default:
    // Continue working
}
```

---

## Part 2: Understanding Channels

### What are Channels?

Channels are Go's pipes for goroutines to communicate. They enable safe concurrent programming without locks.

**Key Concepts:**
- **Send**: `ch <- value` (blocks if channel is full)
- **Receive**: `value := <-ch` (blocks if channel is empty)
- **Close**: `close(ch)` (signals no more values will be sent)
- **Buffered**: `make(chan int, 10)` (holds 10 values before blocking)
- **Unbuffered**: `make(chan int)` (blocks immediately until received)

### Common Patterns

**1. Worker Pool** - Process tasks concurrently with limited workers
```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        go func() {
            for job := range jobs {
                results <- processJob(job)
            }
        }()
    }
}
```

**2. Fan-Out/Fan-In** - Distribute work across multiple goroutines, collect results
```go
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        ch := make(chan int)
        channels[i] = ch
        go worker(input, ch)
    }
    return channels
}
```

**3. Pub/Sub** - Broadcast events to multiple listeners
```go
type PubSub struct {
    subscribers []chan Event
    mu          sync.Mutex
}

func (ps *PubSub) Publish(event Event) {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    for _, ch := range ps.subscribers {
        select {
        case ch <- event:
        default: // Don't block if subscriber is slow
        }
    }
}
```

---

## Part 3: Where to Apply in Your Codebase

### 1. Database Queries (High Priority)

**Current**: Direct DB calls without timeout
**Improvement**: Add context with timeouts

**Files to improve:**
- `internal/sqlite/*-sql.go` (all database queries)
- `internal/dao/dao.go` (generic DAO already has context params!)

**Benefits:**
- Queries auto-cancel when user disconnects
- Prevent slow queries from blocking server
- Better error messages (timeout vs generic error)

---

### 2. Image Processing (Medium Priority)

**Current**: Synchronous image uploads block request
**Improvement**: Use worker pool with channels

**Files to improve:**
- `internal/http/handlers/post-handlers.go` (image uploads)

**Benefits:**
- Handle multiple uploads concurrently
- Limit concurrent processing (prevent memory spikes)
- Return response immediately, process async

---

### 3. Search (Medium Priority)

**Current**: Sequential search across multiple tables
**Improvement**: Fan-out pattern to search concurrently

**Files to improve:**
- `internal/http/handlers/search-handlers.go`

**Benefits:**
- Search users, posts, channels simultaneously
- Faster results (parallelization)
- Better timeout handling

---

### 4. Notification System (Future Enhancement)

**Improvement**: Pub/Sub pattern for real-time notifications

**Benefits:**
- Send notifications without blocking request
- Support multiple notification channels (websocket, email, etc.)

---

## Part 4: Exercises

Complete these exercises in order. Each builds on the previous.

### Exercise 1: Add Context to Database Queries (Beginner)

**Goal**: Learn context-aware database methods

**Task**: Update the `CookieModel.GetUserByCookie` method to use context

**File**: `internal/sqlite/cookies-sql.go:49`

**Instructions**:
1. Add `ctx context.Context` as the first parameter
2. Replace `m.DB.Query` with `m.DB.QueryContext(ctx, ...)`
3. Test with a 5-second timeout from the handler

**Hint**: Context should always be the first parameter

---

### Exercise 2: Database Query Timeout (Beginner)

**Goal**: Prevent slow queries from blocking

**Task**: Create a middleware that adds a 10-second timeout to all database operations

**File**: Create `internal/http/middleware/timeout.go`

**Requirements**:
- Create `WithTimeout(next http.Handler, timeout time.Duration) http.Handler`
- Add timeout context to request context
- All handlers should use this timeout context for DB calls

**Test**: Simulate a slow query with `SELECT * FROM Users WHERE ID IN (SELECT sleep(15))`

---

### Exercise 3: Image Upload Worker Pool (Intermediate)

**Goal**: Process multiple image uploads concurrently

**Task**: Create a worker pool for image processing

**File**: Create `internal/workers/image_worker.go`

**Requirements**:
- Create a buffered channel for image jobs: `make(chan ImageJob, 100)`
- Start 5 worker goroutines that process jobs
- Return job ID immediately to user
- Process image in background
- Store result in database when done

**Architecture**:
```go
type ImageJob struct {
    ID       string
    FilePath string
    UserID   models.UUIDField
}

type ImageWorkerPool struct {
    jobs    chan ImageJob
    workers int
}

func (pool *ImageWorkerPool) Start()
func (pool *ImageWorkerPool) Submit(job ImageJob) error
func (pool *ImageWorkerPool) Shutdown(ctx context.Context) error
```

---

### Exercise 4: Concurrent Search (Intermediate)

**Goal**: Search multiple tables in parallel

**Task**: Refactor search to use fan-out pattern

**File**: `internal/http/handlers/search-handlers.go:20`

**Requirements**:
- Create 3 goroutines: search users, search posts, search channels
- Use channels to collect results
- Use `sync.WaitGroup` or context timeout to wait for all
- Combine results and return

**Pattern**:
```go
type SearchResult struct {
    Users    []models.User
    Posts    []models.Post
    Channels []models.Channel
}

func concurrentSearch(ctx context.Context, query string) (*SearchResult, error) {
    // Your implementation here
}
```

---

### Exercise 5: Graceful Database Shutdown (Advanced)

**Goal**: Close database connections cleanly during shutdown

**Task**: Update server shutdown to wait for in-flight DB queries

**File**: `cmd/server/main.go:108`

**Requirements**:
- Track active database operations using a counter or WaitGroup
- During shutdown, stop accepting new requests
- Wait for existing queries to complete (with timeout)
- Close database connections

**Hint**: Use a done channel to signal shutdown

---

### Exercise 6: Request Tracing with Context (Advanced)

**Goal**: Track requests across handlers for debugging

**Task**: Add request ID to context and log it in all operations

**File**: Create `internal/http/middleware/tracing.go`

**Requirements**:
- Generate unique request ID (UUID)
- Add to context using `context.WithValue`
- Update log statements to include request ID
- Add request ID to response headers

**Bonus**: Track request duration and log slow requests (>1s)

---

### Exercise 7: Circuit Breaker for Database (Expert)

**Goal**: Prevent cascading failures when database is slow

**Task**: Implement a circuit breaker using channels

**File**: Create `internal/circuitbreaker/breaker.go`

**Requirements**:
- Track failure rate in a sliding window
- Open circuit (reject requests) if failure rate > 50%
- Half-open after timeout (allow 1 request to test)
- Close circuit if request succeeds

**States**: Closed (normal) → Open (rejecting) → Half-Open (testing) → Closed

---

## Part 5: Testing Your Improvements

### Benchmark Database Operations

```go
func BenchmarkGetUserWithContext(b *testing.B) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    for i := 0; i < b.N; i++ {
        _, err := db.GetUserByIDContext(ctx, testUserID)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Test Cancellation

```go
func TestQueryCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        time.Sleep(100 * time.Millisecond)
        cancel() // Cancel while query is running
    }()

    _, err := db.QueryContext(ctx, "SELECT sleep(10)")
    if err != context.Canceled {
        t.Errorf("expected context.Canceled, got %v", err)
    }
}
```

---

## Part 6: Performance Gains You Can Expect

**Database Context Timeouts**:
- 20-30% reduction in P99 latency (slow queries don't block)
- Better resource utilization (cancelled queries free DB connections)

**Image Worker Pool**:
- 3-5x throughput for concurrent uploads
- Consistent memory usage (controlled concurrency)

**Concurrent Search**:
- 2-3x faster search results (parallel execution)
- Better user experience (faster response)

---

## Resources

- [Go Context Package](https://pkg.go.dev/context)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency](https://go.dev/blog/context)
- [Database/SQL Context](https://pkg.go.dev/database/sql)

---

## Next Steps

1. Start with Exercise 1 (easiest)
2. Complete exercises 1-3 for immediate performance gains
3. Exercises 4-7 are more advanced - tackle when comfortable
4. Run benchmarks to measure improvements
5. Monitor production metrics (latency, error rates)

Good luck! 🚀

---

## Exercise 4: Async Logging System

### Goal
Build a production-ready async logging system that stores request metrics, errors, and system health data in a database without slowing down your application.

### What You'll Learn
- Applying worker pool pattern to a different domain (logging vs images)
- Non-blocking operations (dropped logs better than slow requests!)
- Database design for analytics and monitoring
- Building admin dashboards with queryable metrics

### Architecture

```
HTTP Request → Middleware → Submit to Queue → Worker Pool → Database
     ↓              ↓                              ↓
  Immediate     Capture                      Async Write
  Response      Metrics                   (Non-blocking!)
```

**Key Insight**: Logging must NEVER slow down requests! We use async workers to decouple logging from request handling.

### Files Created

**Database:**
- `migrations/006_logging_system.sql` - Three tables: RequestLogs, ErrorLogs, SystemMetrics

**Models:**
- `internal/models/logging-models.go` - RequestLog, ErrorLog, SystemMetric structs

**Database Layer:**
- `internal/sqlite/logging-sql.go` - Insert methods, stats queries, cleanup

**Worker Pool:**
- `internal/workers/logger_worker.go` - Async logger pool (YOUR TASK!)
- `internal/workers/logger_worker_test.go` - Test suite

**Middleware:**
- `internal/http/middleware/logging_enhanced.go` - Enhanced logging middleware

### Your Tasks

#### Part 1: Implement Submit Method

Location: `internal/workers/logger_worker.go:86`

**Pattern** (identical to ImageWorkerPool.Submit):
```go
func (pool *LoggerPool) Submit(entry LogEntry) error {
    // 1. Check shutdown
    if pool.isShutdown.Load() {
        return fmt.Errorf("logger pool is shut down")
    }

    // 2. Non-blocking send
    select {
    case pool.logs <- entry:
        return nil
    default:
        // IMPORTANT: Increment dropped counter!
        pool.droppedLogs.Add(1)
        return fmt.Errorf("logger pool queue is full (dropped: %d)", pool.droppedLogs.Load())
    }
}
```

**Why dropped logs counter?** Admin dashboard can show "X logs dropped in last hour" - indicates queue too small or DB too slow!

#### Part 2: Implement Shutdown Method

Location: `internal/workers/logger_worker.go:105`

**Pattern** (identical to ImageWorkerPool.Shutdown):
```go
func (pool *LoggerPool) Shutdown(ctx context.Context) error {
    pool.isShutdown.Store(true)
    close(pool.shutdownCh)
    
    done := make(chan struct{})
    go func() {
        pool.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

#### Part 3: Extract User ID from Context

Location: `internal/http/middleware/logging_enhanced.go:47`

**Pattern:**
```go
// Get user from context (set by WithUser middleware)
user, ok := r.Context().Value("user").(*models.User)

var userID models.UUIDField
if ok && user != nil {
    userID = user.ID
} else {
    userID = models.ZeroUUIDField()  // Nil UUID for anonymous requests
}
```

### Testing

```sh
# Run logger worker tests
go test -v ./internal/workers/... -run TestLogger

# Run database migration
bin/codex migrate

# Expected output:
# ✓ Successfully applied migration: migrations/006_logging_system.sql
```

### Integration

Once implemented, you'll need to:

1. **Initialize logger pool in `cmd/server/main.go`:**
```go
// Create logger pool
loggerPool := workers.NewLoggerPool(3, 1000, db)
loggerPool.Start()

// Graceful shutdown
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    loggerPool.Shutdown(ctx)
}()
```

2. **Update middleware in `internal/http/routes/routes.go`:**
```go
// Replace old logging middleware
// mux.Use(middleware.Logging)

// With enhanced logging
mux.Use(middleware.LoggingEnhanced(loggerPool))
```

### Admin Dashboard Queries

After implementing, you can build an admin page with queries like:

```go
// Get request stats for last 24 hours
stats, _ := loggingModel.GetRequestStats("datetime('now', '-24 hours')")
fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
fmt.Printf("Avg Duration: %.2fms\n", stats.AvgDuration)
fmt.Printf("Error Rate: %.2f%%\n", stats.ErrorRate)
fmt.Printf("Unique Users: %d\n", stats.UniqueUsers)

// Get slowest endpoints
SELECT Path, AVG(Duration) as avg_duration
FROM RequestLogs
WHERE Timestamp >= datetime('now', '-1 hour')
GROUP BY Path
ORDER BY avg_duration DESC
LIMIT 10;

// Get error patterns
SELECT COUNT(*) as count, Message
FROM ErrorLogs
WHERE Timestamp >= datetime('now', '-1 day')
GROUP BY Message
ORDER BY count DESC
LIMIT 20;
```

### Key Differences from Image Worker

| Aspect | Image Worker | Logger Worker |
|--------|--------------|---------------|
| **Critical Path** | Can wait briefly (user uploads) | Must be instant (every request!) |
| **Queue Full Strategy** | Return 503 error to user | **Drop log** (don't slow request) |
| **Failure Handling** | Retry, user sees error | Log to console, continue |
| **Performance Impact** | Only on uploads | **Every single request** |
| **Monitoring** | Jobs processed | **Dropped logs counter** |

### Success Criteria

✅ All tests pass  
✅ Migration applies cleanly  
✅ Submit() is non-blocking  
✅ Shutdown() waits for pending logs  
✅ Dropped logs counter increments when queue full  
✅ No database writes block HTTP responses  

### Next Steps

After completing this exercise, you can:
1. Build an admin dashboard (HTML + charts)
2. Add alerting (email admin if error rate > 5%)
3. Auto-cleanup old logs (cron job using CleanupOldLogs)
4. Export metrics to monitoring tools (Prometheus, Grafana)

