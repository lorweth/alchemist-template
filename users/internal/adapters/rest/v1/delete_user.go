package v1

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
)

func (hdl UserHandler) DeleteUser() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// Get user ID from URL
		userID, err := httpio.URLParam[int64](r, "id")
		if err != nil {
			return err
		}

		// Validate user ID
		if userID == 0 {
			return errUserIDMustBeGreaterThanZero
		}

		// Delete user
		if err := hdl.svc.Delete(ctx, domain.User{
			ID: userID,
		}); err != nil {
			return convertServiceError(err)
		}

		// Write response
		httpio.WriteJSON(w, r, httpio.Response[httpio.Message]{
			Status: http.StatusOK,
			Body: httpio.Message{
				Code: "delete_success",
				Desc: "Delete successfully",
			},
		})

		return nil
	})
}
