package services

import (
	"errors"
)

var (
	EmailHasBeenUsed = errors.New("email_has_been_used")
	UserNotFound     = errors.New("user_not_found")
)
