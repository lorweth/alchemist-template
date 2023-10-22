package controller

import (
	"context"
	"database/sql"
	"errors"

	"github.com/virsavik/alchemist-template/users/internal/model"
	"github.com/virsavik/alchemist-template/users/internal/repository"
)

type userController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) UserController {
	return &userController{
		repo: repo,
	}
}

func (c userController) RegisterUser(ctx context.Context, user *model.User) error {
	existed, err := c.repo.FindOne(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {

			return err
		}
	}

	if existed != nil && existed.ID != 0 {
		return ErrEmailAlreadyInUse
	}

	if err := c.repo.Save(ctx, user); err != nil {
		return err
	}

	return nil
}

func (c userController) GetUser(ctx context.Context, email string) (*model.User, error) {
	user, err := c.repo.FindOne(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c userController) EnableUser(ctx context.Context, email string) error {
	existed, err := c.repo.FindOne(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserDoesNotExist
		}

		return nil
	}

	if existed == nil {
		return ErrUserDoesNotExist
	}

	existed.IsActive = true
	if err := c.repo.Save(ctx, existed); err != nil {
		return err
	}

	return nil
}

func (c userController) DisableUser(ctx context.Context, email string) error {
	existed, err := c.repo.FindOne(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserDoesNotExist
		}

		return nil
	}

	if existed == nil {
		return ErrUserDoesNotExist
	}

	existed.IsActive = false
	if err := c.repo.Save(ctx, existed); err != nil {
		return err
	}

	return nil
}
