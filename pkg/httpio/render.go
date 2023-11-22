package httpio

import (
	"encoding/json"
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/logger"
)

type Response[T any] struct {
	Status  int
	Headers map[string]string
	Body    T
}

func WriteJSON[T any](w http.ResponseWriter, r *http.Request, data Response[T]) {
	ctx := r.Context()

	// Update response headers
	for header, val := range data.Headers {
		w.Header().Set(header, val)
	}

	w.Header().Set("Content-Type", "application/json")

	// Update status code
	w.WriteHeader(data.Status)

	if err := json.NewEncoder(w).Encode(data.Body); err != nil {
		logger.FromCtx(ctx).Errorf(err, "json encode error")
	}
}
