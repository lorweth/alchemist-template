package pagination

import "errors"

var (
	ErrPageNumberInvalid = errors.New("page number is invalid")
	ErrPageSizeInvalid   = errors.New("page size is invalid")
)
