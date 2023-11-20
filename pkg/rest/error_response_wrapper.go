package rest

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

// ErrResponseWrapper write internal server error to client when got unexpected error
func ErrResponseWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		if err := fn(w, r); err != nil {
			var apiErr respond.Error
			if errors.As(err, &apiErr) {
				respond.FromError(apiErr).WriteJSON(ctx, w)

				return
			}

			// Internal Server error
			respond.InternalServerError(respond.Message{Name: "internal_server_error"}).WriteJSON(ctx, w)
			span.RecordError(err, trace.WithStackTrace(true))
		}
	}
}
