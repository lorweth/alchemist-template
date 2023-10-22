package controller

import "errors"

var (
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrUserDoesNotExist  = errors.New("user does not exist")
)
