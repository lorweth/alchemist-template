package controller

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/model"
)

type UserController interface {
	RegisterUser(ctx context.Context, user *model.User) error

	GetUser(ctx context.Context, email string) (*model.User, error)

	EnableUser(ctx context.Context, email string) error

	DisableUser(ctx context.Context, email string) error
}
