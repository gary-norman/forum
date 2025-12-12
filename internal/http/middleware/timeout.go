package middleware

import (
	"context"
	"net/http"
	"time"
)

// TODO(human): Exercise 2 - Create timeout middleware
//
// Instructions:
// 1. Create a function called WithTimeout that takes:
//    - next http.Handler (the next handler in the chain)
//    - timeout time.Duration (how long to wait before timing out)
//    Returns: http.Handler
//
// 2. Inside WithTimeout, return an http.HandlerFunc that:
//    a) Creates a context with timeout: ctx, cancel := context.WithTimeout(r.Context(), timeout)
//    b) Always defers cancel() to prevent context leaks: defer cancel()
//    c) Passes the request with the new context to next handler: next.ServeHTTP(w, r.WithContext(ctx))
//
// Pattern to follow:
// func WithTimeout(next http.Handler, timeout time.Duration) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Your implementation here
//     })
// }
//
// Hint: This is very similar to the WithUser middleware in auth.go!

// Your implementation here:
func WithTimeout(next http.Handler, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
