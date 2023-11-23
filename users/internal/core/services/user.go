package services

import (
	"context"

	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc UserService) GetAll(ctx context.Context, input ports.GetUserInput) ([]domain.User, int64, error) {
	users, total, err := svc.repo.GetAll(ctx, input)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (svc UserService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	// Find user by email
	selectedUser, err := svc.repo.GetOne(ctx, ports.GetUserInput{
		Email: user.Email,
	})
	if err != nil {
		return domain.User{}, err
	}

	// Return error if user exists
	if selectedUser.ID != 0 {
		return domain.User{}, EmailHasBeenUsed
	}

	// Save user
	createdUser, err := svc.repo.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (svc UserService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	// Find user by email
	selectedUser, err := svc.repo.GetOne(ctx, ports.GetUserInput{
		Email: user.Email,
	})
	if err != nil {
		return domain.User{}, err
	}

	// Return error if user not exists
	if selectedUser.ID == 0 {
		return domain.User{}, UserNotFound
	}

	// Save user
	createdUser, err := svc.repo.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (svc UserService) Delete(ctx context.Context, user domain.User) error {
	// Find user by email
	selectedUser, err := svc.repo.GetOne(ctx, ports.GetUserInput{
		ID: user.ID,
	})
	if err != nil {
		return err
	}

	// Return error if user not exists
	if selectedUser.ID == 0 {
		return UserNotFound
	}

	// Delete user
	if err := svc.repo.Delete(ctx, selectedUser); err != nil {
		return err
	}

	return nil
}
