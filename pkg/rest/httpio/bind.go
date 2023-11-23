package httpio

import (
	"encoding/json"
	"io"
)

// BindJSON decodes JSON data from the given io.Reader into a provided type.
// It returns an instance of the provided type and any decoding error encountered.
//
// Example usage:
//
//	data, err := BindJSON[MyStruct](req.Body)
//	if err != nil {
//	    // Handle error
//	}
//	// Use the decoded data
func BindJSON[T any](reqBody io.Reader) (T, error) {
	var req T
	if err := json.NewDecoder(reqBody).Decode(&req); err != nil {
		return req, err
	}

	return req, nil
}
