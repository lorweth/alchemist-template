package repository

import (
	"errors"
)

var (
	ErrUserIDInvalid     = errors.New("user id is invalid")
	ErrMoreThanOneRecord = errors.New("more than one record match with given input")
)
