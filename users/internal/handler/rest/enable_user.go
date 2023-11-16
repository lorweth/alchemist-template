package rest

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest"
	"github.com/virsavik/alchemist-template/pkg/rest/request"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

type EnableUserRequest struct {
	Email string `json:"email"`
}

type EnableUserResponse struct {
	Success bool `json:"success"`
}

// EnableUser set user status to active
func (h Handler) EnableUser() http.HandlerFunc {
	return rest.ErrResponseWrapper(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := request.BindJSON[EnableUserRequest](r.Body)
		if err != nil {
			return err
		}

		if err := h.userCtrl().EnableUser(r.Context(), req.Email); err != nil {
			return convertCtrlErr(err)
		}

		respond.Ok(EnableUserResponse{
			Success: true,
		}).WriteJSON(ctx, w)

		return nil
	})
}
