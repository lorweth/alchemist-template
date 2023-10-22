package repository

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/model"
)

type UserRepository interface {
	FindOne(ctx context.Context, email string) (*model.User, error)

	Save(ctx context.Context, user *model.User) error

	Delete(ctx context.Context, email string) error
}
