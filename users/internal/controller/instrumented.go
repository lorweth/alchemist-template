package controller

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/virsavik/alchemist-template/users/internal/model"
)

type instrumentedUserController struct {
	UserController

	userRegistered prometheus.Counter
}

func New(ctrl UserController, userRegistered prometheus.Counter) UserController {
	return instrumentedUserController{
		UserController: ctrl,
		userRegistered: userRegistered,
	}
}

func (i instrumentedUserController) RegisterUser(ctx context.Context, user *model.User) error {
	if err := i.UserController.RegisterUser(ctx, user); err != nil {
		return err
	}

	i.userRegistered.Inc()
	return nil
}
