package httpio

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	pkgerrors "github.com/pkg/errors"
)

// URLQuery retrieves a query parameter named 'key' from the given HTTP request,
// converts it to the specified type 'T', and returns the result.
// If the query parameter is not present or the conversion fails, the function returns the default value for 'T'.
// It supports types int, int64, string, bool, and struct (which is unmarshaled from JSON).
// The function uses reflection to dynamically handle different types.
func URLQuery[T any](r *http.Request, key string) (T, error) {
	qVal := r.URL.Query().Get(key)

	var rs T

	v := reflect.ValueOf(&rs).Elem()
	if !v.CanAddr() || !v.CanSet() {
		return rs, errors.New("cannot assign value")
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		// Assuming val is a string representing an integer
		intVal, err := strconv.Atoi(qVal)
		if err != nil {
			// Ignore error & use default value
		}

		v.SetInt(int64(intVal))

	case reflect.String:
		// Assuming val is already a string
		v.SetString(qVal)

	case reflect.Bool:
		boolV, err := strconv.ParseBool(qVal)
		if err != nil {
			// Ignore err & user default value
		}
		v.SetBool(boolV)

	case reflect.Struct:
		if err := json.Unmarshal([]byte(qVal), &rs); err != nil {
			return rs, pkgerrors.WithStack(err)
		}
	}

	return rs, nil
}
