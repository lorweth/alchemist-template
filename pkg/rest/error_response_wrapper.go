package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/logger"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

// ErrResponseWrapper write internal server error to client when got unexpected error
func ErrResponseWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("content-type", "application/json")

		if err := fn(w, r); err != nil {
			var apiErr respond.Error
			if errors.As(err, &apiErr) {
				w.WriteHeader(apiErr.Code)
				recordJSONError(ctx, json.NewEncoder(w).Encode(apiErr))

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			recordJSONError(ctx, json.NewEncoder(w).Encode(respond.ErrInternalServer))
		}
	}
}

func recordJSONError(ctx context.Context, err error) {
	if err != nil {
		logger.FromCtx(ctx).Error(err, "error when encode json")

		trace.SpanFromContext(ctx).AddEvent("Encode", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
	}
}
