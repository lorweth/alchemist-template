package httpio

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// HandlerFunc wraps an HTTP handler function that returns an error.
// It adds OpenTelemetry tracing and handles specific error types by responding with JSON.
// If the error is of type httpio.Error, a custom JSON response is generated.
// If the error is not of type httpio.Error, a generic internal server error response is generated.
func HandlerFunc(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		// Execute the provided handler function and handle errors
		if err := fn(w, r); err != nil {
			var apiErr Error
			if errors.As(err, &apiErr) {
				WriteJSON(w, r, Response[Message]{
					Status: apiErr.Status,
					Body: Message{
						Key:     apiErr.Key,
						Content: apiErr.Desc,
					},
				})

				return
			}

			// If the error is not of type "Error", respond with a generic internal server error
			WriteJSON(w, r, Response[Message]{
				Status: http.StatusInternalServerError,
				Body:   MsgInternalServerError,
			})

			// Record the error in the OpenTelemetry span with a stack trace
			span.RecordError(err, trace.WithStackTrace(true))
		}
	}
}
