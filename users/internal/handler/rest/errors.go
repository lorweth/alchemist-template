package rest

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest/respond"
	"github.com/virsavik/alchemist-template/users/internal/controller"
)

func convertCtrlErr(err error) error {
	switch err.Error() {
	case controller.ErrEmailAlreadyInUse.Error(),
		controller.ErrUserDoesNotExist.Error():
		return respond.Error{
			Code: http.StatusBadRequest,
			Desc: err.Error(),
		}
	default:
		return err
	}
}
