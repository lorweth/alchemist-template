package v1

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/services"
)

var (
	errEmailCannotBeBlank          = httpio.Error{Status: http.StatusBadRequest, Code: "email_blank", Desc: "Email cannot be blank"}
	errUserIDMustBeGreaterThanZero = httpio.Error{Status: http.StatusBadRequest, Code: "user_id_zero", Desc: "User id must be greater than zero"}
	errToTimeCannotBeforeFromTime  = httpio.Error{Status: http.StatusBadRequest, Code: "from_time_cannot_before_to_time", Desc: "From time cannot before to time"}
)

func convertServiceError(err error) error {
	switch err.Error() {
	case services.EmailHasBeenUsed.Error(),
		services.UserNotFound.Error():
		return httpio.Error{
			Status: http.StatusBadRequest,
			Code:   err.Error(),
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
		Code:   errName,
		Desc:   err.Error(),
	}
}
