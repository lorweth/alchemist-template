package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Otel is a middleware function that provides OpenTelemetry (OTel) tracing capabilities
// for an HTTP server. It takes a trace.Tracer as a parameter, which allows you to
// use an OTel-compatible tracer for tracing, and returns an HTTP middleware handler.
//
// The returned middleware sets up a new span within the provided tracer for each incoming
// HTTP request. It logs tracing information related to the request, including the HTTP route,
// HTTP method, and other attributes, and sets the span's kind as 'Server'.
//
// It also wraps the original HTTP handler to capture the response's status code, allowing
// it to be added as an attribute to the span. This facilitates tracing and monitoring of
// requests and responses.
//
// Example usage:
//
//	tracerProvider := otel.NewTracerProvider()
//	tracer := tracerProvider.Tracer("my-service-name")
//	http.Handle("/", Otel(tracer)(myHandler))
//	http.ListenAndServe(":8080", nil)
func Otel(tracer trace.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(
				r.Context(),
				"chi-route",
				trace.WithAttributes(
					semconv.HTTPRouteKey.String(r.URL.Path),
					semconv.HTTPMethod(r.Method),
				),
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// Wrap the response writer to get status_code easily
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r.WithContext(ctx))

			// Add the status code as an attribute to the span.
			span.SetAttributes(attribute.Int("http.status_code", ww.Status()))
		}

		return http.HandlerFunc(fn)
	}
}
