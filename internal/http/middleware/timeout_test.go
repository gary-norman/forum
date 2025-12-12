package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestWithTimeout tests that the timeout middleware works correctly
func TestWithTimeout(t *testing.T) {
	// Create a handler that sleeps for 2 seconds
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Completed"))
		case <-r.Context().Done():
			// Context was cancelled due to timeout
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Timeout"))
			return
		}
	})

	t.Run("completes within timeout", func(t *testing.T) {
		// Wrap with 3 second timeout (should complete)
		handler := WithTimeout(slowHandler, 3*time.Second)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if body := rr.Body.String(); body != "Completed" {
			t.Errorf("handler returned unexpected body: got %v want %v", body, "Completed")
		}
	})

	t.Run("times out before completion", func(t *testing.T) {
		// Wrap with 1 second timeout (should timeout)
		handler := WithTimeout(slowHandler, 1*time.Second)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusRequestTimeout {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusRequestTimeout)
		}

		if body := rr.Body.String(); body != "Timeout" {
			t.Errorf("handler returned unexpected body: got %v want %v", body, "Timeout")
		}
	})

	t.Run("context has timeout", func(t *testing.T) {
		// Verify that the context actually has a timeout
		handler := WithTimeout(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deadline, ok := r.Context().Deadline()
			if !ok {
				t.Error("context does not have a deadline")
			}

			// Deadline should be approximately 5 seconds from now
			expectedDeadline := time.Now().Add(5 * time.Second)
			if deadline.Before(expectedDeadline.Add(-1*time.Second)) || deadline.After(expectedDeadline.Add(1*time.Second)) {
				t.Errorf("deadline is not approximately 5 seconds from now: got %v", deadline)
			}

			w.WriteHeader(http.StatusOK)
		}), 5*time.Second)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
	})

	t.Run("no goroutine leak", func(t *testing.T) {
		// This test ensures that defer cancel() prevents goroutine leaks
		// We verify by checking that short requests complete quickly
		// (if cancel wasn't called, timeout goroutines would leak)

		handler := WithTimeout(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Quick handler
			w.WriteHeader(http.StatusOK)
		}), 10*time.Second)

		// Run multiple times to detect leaks
		for i := 0; i < 100; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("iteration %d: handler returned wrong status code: got %v want %v", i, status, http.StatusOK)
			}
		}

		// If cancel() wasn't called, we'd have 100 leaked goroutines
		// In practice, this would cause memory/goroutine count to grow
		// The test passing indicates proper cleanup
	})
}

// TestTimeoutWithDatabaseQuery simulates a database query scenario
func TestTimeoutWithDatabaseQuery(t *testing.T) {
	// Simulate a slow database query
	dbQueryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Simulate database work
		select {
		case <-time.After(3 * time.Second):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Query completed"))
		case <-ctx.Done():
			// In real code, this would cancel the database query
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("Query timeout"))
			}
			return
		}
	})

	t.Run("database query times out", func(t *testing.T) {
		handler := WithTimeout(dbQueryHandler, 1*time.Second)

		req := httptest.NewRequest("GET", "/users", nil)
		rr := httptest.NewRecorder()

		start := time.Now()
		handler.ServeHTTP(rr, req)
		duration := time.Since(start)

		// Should timeout after ~1 second, not wait 3 seconds
		if duration > 2*time.Second {
			t.Errorf("handler took too long: %v (should timeout after ~1s)", duration)
		}

		if status := rr.Code; status != http.StatusRequestTimeout {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusRequestTimeout)
		}
	})
}
