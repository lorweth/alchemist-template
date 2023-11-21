package services

import (
	"errors"
)

var (
	EmailAlreadyLinkedToAnAccount = errors.New("email already linked to an account")
)
