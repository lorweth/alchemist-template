package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/logger"
)

// Logger is a middleware function that provides request logging capabilities
// for an HTTP server. It takes a logger.Logger as a parameter, which allows you
// to use a custom logger for logging, and returns an HTTP middleware handler.
// If a custom logger is not provided (i.e., l is nil), it uses the default
// middleware.Logger for logging.
//
// The returned middleware logs information about incoming HTTP requests, such as
// the request protocol, path, request ID, response status, and response size.
// This information is included in log entries to track and monitor server activity.
//
// The middleware wraps the original HTTP handler and logs the details of each
// request and response. It also sets a new context with the provided logger to
// allow other parts of the application to access and use the same logger instance.
//
// Example usage:
//
//	customLogger := myCustomLoggerImplementation()
//	http.Handle("/", Logger(customLogger)(myHandler))
//	http.ListenAndServe(":8080", nil)
func Logger(l logger.Logger) func(next http.Handler) http.Handler {
	// Use default logger when custom logger not provided
	if l == nil {
		return middleware.Logger
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				reqLogger := l.With(
					logger.WithString("host.name", r.Host),
					logger.WithString("url.path", r.URL.Path),
					logger.WithString("url.query", r.URL.RawQuery),
					logger.WithString("http.request.method_original", r.Method),
					logger.WithInt("http.request.body.size", int(r.ContentLength)),
					logger.WithString("http.request.proto", r.Proto),
					logger.WithString("http.request.remote_address", r.RemoteAddr),
					logger.WithString("user_agent.original", r.UserAgent()),
					logger.WithInt("http.response.status_code", ww.Status()),
					logger.WithInt("http.response.body.size", ww.BytesWritten()),
					logger.WithString("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()),
				)

				reqLogger.Infof("Served")
			}()

			next.ServeHTTP(ww, r.WithContext(logger.SetInCtx(r.Context(), l)))
		}

		return http.HandlerFunc(fn)
	}
}
