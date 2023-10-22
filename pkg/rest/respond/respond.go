package respond

import (
	"encoding/json"
	"net/http"
)

// Response representing for http response data
type Response[T any] struct {
	status int
	data   T
}

// WriteJSON encode the JSON of Response data to the stream, followed by a newline character
func (resp Response[T]) WriteJSON(w http.ResponseWriter) error {
	w.WriteHeader(resp.status)
	return json.NewEncoder(w).Encode(resp.data)
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
