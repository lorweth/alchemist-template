package httpio

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

func Handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		if err := fn(w, r); err != nil {
			var apiErr respond.Error
			if errors.As(err, &apiErr) {
				WriteJSON(w, r, Response[M]{
					Status: apiErr.Status,
					Body: M{
						"key":     apiErr.Key,
						"message": apiErr.Message,
					},
				})

				return
			}

			// Internal Server error
			WriteJSON(w, r, Response[M]{
				Status: http.StatusInternalServerError,
				Body: M{
					"key":     "internal_server_error",
					"message": "Internal server error",
				},
			})
			span.RecordError(err, trace.WithStackTrace(true))
		}
	}
}
