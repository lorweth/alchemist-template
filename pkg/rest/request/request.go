package request

import (
	"encoding/json"
	"io"
)

// BindJSON deserialize JSON from request
func BindJSON[T any](r io.Reader) (T, error) {
	var req T
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return req, err
	}

	return req, nil
}
