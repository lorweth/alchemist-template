package v1

import (
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

type UserHandler struct {
	svc ports.UserService
}

func NewUserHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}
