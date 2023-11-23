package httpio

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// URLParam extracts a URL parameter named 'key' from the given HTTP request and converts it to the specified type 'T'.
// The function supports types int, int64, and string.
// Note: The provided type 'T' should be a concrete type (not a pointer).
func URLParam[T int | int64 | string](r *http.Request, key string) (T, error) {
	// Get URL param value
	str := chi.URLParam(r, key)

	var rs T

	v := reflect.ValueOf(&rs).Elem()
	if !v.CanAddr() || !v.CanSet() {
		return rs, errors.New("cannot assign value")
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		// Assuming val is a string representing an integer
		intVal, err := strconv.Atoi(str)
		if err != nil {
			// Ignore error & use default value
		}

		v.SetInt(int64(intVal))

	case reflect.String:
		// Assuming val is already a string
		v.SetString(str)
	}

	return rs, nil
}
