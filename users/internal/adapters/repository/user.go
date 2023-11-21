package repository

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

func (r Repository) GetAll(ctx context.Context, input ports.GetUserInput) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Save(ctx context.Context, user domain.User) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

