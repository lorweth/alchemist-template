package v1

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/services"
)

var (
	errEmailCannotBeBlank          = httpio.Error{Status: http.StatusBadRequest, Key: "email_blank", Desc: "Email cannot be blank"}
	errUserIDMustBeGreaterThanZero = httpio.Error{Status: http.StatusBadRequest, Key: "user_id_zero", Desc: "User id must be greater than zero"}
)

func convertServiceError(err error) error {
	switch err.Error() {
	case services.EmailHasBeenUsed.Error(),
		services.UserNotFound.Error():
		return httpio.Error{
			Status: http.StatusBadRequest,
			Key:    err.Error(),
		}
	default:
		return err
	}
}

func wrapBadRequestError(err error, name string) error {
	errName := "invalid_request"
	if name != "" {
		errName = name
	}

	return httpio.Error{
		Status: http.StatusBadRequest,
		Key:    errName,
		Desc:   err.Error(),
	}
}
