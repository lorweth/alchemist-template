package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	scopeName = "github.com/virsavik/alchemist-template/pkg/rest/middleware/tracer"
	version   = "1.0.0"
)

// OtelTracer is a middleware function that provides OpenTelemetry (OTel) tracing capabilities
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
//	http.Handle("/", OtelTracer(tracer)(myHandler))
//	http.ListenAndServe(":8080", nil)
func OtelTracer() func(next http.Handler) http.Handler {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer(scopeName,
		trace.WithInstrumentationVersion(version),
	)
	propagators := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := tracer.Start(
				ctx,
				"http.request",
				trace.WithAttributes(
					semconv.HostName(r.Host),
					semconv.URLPath(r.URL.Path),
					semconv.URLQuery(r.URL.RawQuery),
					semconv.HTTPRequestMethodOriginal(r.Method),
					semconv.HTTPRequestBodySize(int(r.ContentLength)),
					attribute.String("http.request.proto", r.Proto),
					attribute.String("http.request.remote_address", r.RemoteAddr),
					semconv.UserAgentOriginal(r.UserAgent()),
				),
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// Wrap the response writer to get status_code easily
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r.WithContext(ctx))

			// Add the status code as an attribute to the span.
			span.SetAttributes(
				semconv.HTTPResponseStatusCode(ww.Status()),
				semconv.HTTPResponseBodySize(ww.BytesWritten()),
			)
		}

		return http.HandlerFunc(fn)
	}
}
