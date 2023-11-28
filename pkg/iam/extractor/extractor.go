package extractor

import (
	"net/http"
	"strings"
)

// AuthHeader is a token extractor that takes a request
// and extracts the token from the Authorization header.
func AuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no JWT.
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", ErrInvalidFormat
	}

	return authHeaderParts[1], nil
}

// Parameter returns an Extractor that extracts the token
// from the specified query string parameter.
func Parameter(param string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		return r.URL.Query().Get(param), nil
	}
}
