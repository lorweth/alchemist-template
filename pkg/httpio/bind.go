package httpio

import (
	"encoding/json"
	"io"
)

func BindJSON[T any](reqBody io.Reader) (T, error) {
	var req T
	if err := json.NewDecoder(reqBody).Decode(&req); err != nil {
		return req, err
	}

	return req, nil
}
