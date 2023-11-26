package iam

import (
	"errors"
)

var (
	ErrTokenSubNotFound = errors.New("token sub missing")
	ErrTokenInvalid     = errors.New("token invalid")
)
