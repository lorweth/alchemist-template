package respond

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/logger"
)

func New[T any](status int, data T) Response[T] {
	return Response[T]{
		status: status,
		data:   data,
	}
}

func FromError(err Error) Response[Message] {
	return Response[Message]{
		status: err.StatusCode,
		data: Message{
			Name:    err.Name,
			Message: err.Message,
		},
	}
}

// WriteJSON encode the JSON of Response data
func (resp Response[T]) WriteJSON(ctx context.Context, w http.ResponseWriter) {
	resp.WriteJSONWithHeader(ctx, w, nil)
}

// WriteJSONWithHeader encode the JSON of Response data with header
func (resp Response[T]) WriteJSONWithHeader(ctx context.Context, w http.ResponseWriter, header map[string]string) {
	w.Header().Set("content-type", "application/json")
	for key, val := range header {
		w.Header().Set(key, val)
	}

	w.WriteHeader(resp.status)

	if err := json.NewEncoder(w).Encode(resp.data); err != nil {
		logger.FromCtx(ctx).Errorf(err, "json encode error")
	}
}

// Ok creates an OK(200) Response
func Ok[T any](data T) Response[T] {
	return Response[T]{
		status: http.StatusOK,
		data:   data,
	}
}

// BadRequest creates an Bad request (400) Response
func BadRequest[T any](data T) Response[T] {
	return Response[T]{
		status: http.StatusBadRequest,
		data:   data,
	}
}

// Unauthorized creates an Unauthorized (401) Response
func Unauthorized[T any](data T) Response[T] {
	return Response[T]{
		status: http.StatusUnauthorized,
		data:   data,
	}
}

func InternalServerError[T any](data T) Response[T] {
	return Response[T]{
		status: http.StatusInternalServerError,
		data:   data,
	}
}
