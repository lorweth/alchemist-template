package ports

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/core/domain"
)

type GetUserInput struct {
	IDs []int64
}

type UserRepository interface {
	GetAll(ctx context.Context, input GetUserInput) ([]domain.User, error)

	GetByEmail(ctx context.Context, email string) (domain.User, error)

	Save(ctx context.Context, user domain.User) (domain.User, error)
}

type UserService interface {
	GetAll(ctx context.Context, input GetUserInput) ([]domain.User, error)

	Create(ctx context.Context, user domain.User) (domain.User, error)
}
