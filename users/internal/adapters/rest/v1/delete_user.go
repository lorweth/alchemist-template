package v1

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest"
	"github.com/virsavik/alchemist-template/pkg/rest/request"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
)

func (hdl UserHandler) DeleteUser() http.HandlerFunc {
	return rest.ErrResponseWrapper(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// Decode request
		req, err := request.BindJSON[deleteUserRequest](r.Body)
		if err != nil {
			return err
		}

		// Validate request
		if err := req.IsValid(); err != nil {
			return err
		}

		// Delete user
		if err := hdl.svc.Delete(ctx, domain.User{
			ID: req.ID,
		}); err != nil {
			return convertServiceError(err)
		}

		// Write response
		respond.Ok(respond.Message{Name: "delete_success", Message: "Delete user successfully"}).WriteJSON(ctx, w)

		return nil
	})
}

type deleteUserRequest struct {
	ID int64
}

func (req deleteUserRequest) IsValid() error {
	if req.ID == 0 {
		return errUserIDMustBeGreaterThanZero
	}

	return nil
}
