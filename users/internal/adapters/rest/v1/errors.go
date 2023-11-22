package v1

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest/respond"
	"github.com/virsavik/alchemist-template/users/internal/core/services"
)

var (
	errEmailCannotBeBlank          = respond.Error{Status: http.StatusBadRequest, Key: "email_blank", Message: "Email cannot be blank"}
	errUserIDMustBeGreaterThanZero = respond.Error{Status: http.StatusBadRequest, Key: "user_id_zero", Message: "User id must be greater than zero"}
)

func convertServiceError(err error) error {
	switch err.Error() {
	case services.EmailHasBeenUsed.Error(),
		services.UserNotFound.Error():
		return respond.Error{
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

	return respond.Error{
		Status:  http.StatusBadRequest,
		Key:     errName,
		Message: err.Error(),
	}
}
