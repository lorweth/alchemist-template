package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestID is a middleware function that generates and manages unique request IDs
// for incoming HTTP requests. It returns an HTTP middleware handler.
//
// The middleware adds a unique request ID to the context of each incoming request,
// which can be used for tracking and correlating requests and responses. This request
// ID is typically included in log entries and can be helpful for debugging and monitoring.
//
// Example usage:
//
//	http.Handle("/", RequestID()(myHandler))
//	http.ListenAndServe(":8080", nil)
func RequestID() func(next http.Handler) http.Handler {
	fn := middleware.RequestID
	return fn
}
