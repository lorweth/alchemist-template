package ports

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/core/domain"
)

type GetUserInput struct {
	ID    int64
	Email string
}

type UserRepository interface {
	GetAll(ctx context.Context, input GetUserInput) ([]domain.User, error)

	GetOne(ctx context.Context, input GetUserInput) (domain.User, error)

	Save(ctx context.Context, user domain.User) (domain.User, error)

	Delete(ctx context.Context, user domain.User) error
}

type UserService interface {
	GetAll(ctx context.Context, input GetUserInput) ([]domain.User, error)

	Create(ctx context.Context, user domain.User) (domain.User, error)

	Update(ctx context.Context, user domain.User) (domain.User, error)

	Delete(ctx context.Context, user domain.User) error
}
