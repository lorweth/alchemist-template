package services

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func (serv UserService) GetAll(ctx context.Context, input ports.GetUserInput) ([]domain.User, error) {
	users, err := serv.repo.GetAll(ctx, input)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (serv UserService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	selectedUser, err := serv.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}

	if selectedUser.ID != 0 {
		return domain.User{}, EmailAlreadyLinkedToAnAccount
	}

	createdUser, err := serv.repo.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}
