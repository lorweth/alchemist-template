package rest

import (
	"github.com/virsavik/alchemist-template/pkg/di"
	"github.com/virsavik/alchemist-template/users/internal/constants"
	"github.com/virsavik/alchemist-template/users/internal/controller"
)

type Handler struct {
	ctn di.Container
}

func New(ctn di.Container) Handler {
	return Handler{
		ctn: ctn,
	}
}

func (h Handler) userCtrl() controller.UserController {
	return h.ctn.Get(constants.UsersCtrlKey).(controller.UserController)
}
