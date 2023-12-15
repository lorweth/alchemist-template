package extractor

import (
	"errors"
)

var (
	ErrInvalidFormat = errors.New("authorization header format must be bearer token")
)
