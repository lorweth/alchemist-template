package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/logger"
)

// Recover is a middleware function that recovers from panics in the application
// It logs the panic details, captures the entire stack trace, and returns a 500 Internal Server Error response.
// This middleware is designed to be used in the HTTP request handling pipeline.
func Recover() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// Recover from the panic and obtain the panic value.
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%+v", p)
					}

					// Capture and log the entire stack trace along with the error details.
					logger.FromCtx(r.Context()).Errorf(err, "caught a panic, stacktrace: %s", debug.Stack())

					// Respond with a 500 Internal Server Error and log any encoding errors.
					httpio.WriteJSON(w, r, httpio.Response[httpio.Message]{
						Status: http.StatusInternalServerError,
						Body:   httpio.MsgInternalServerError,
					})

					// Record error
					span := trace.SpanFromContext(r.Context())
					span.SetAttributes(attribute.Bool("error", true))
					span.RecordError(err, trace.WithStackTrace(true))
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
