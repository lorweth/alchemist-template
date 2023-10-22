package rest

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest"
	"github.com/virsavik/alchemist-template/pkg/rest/request"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

type DisableUserRequest struct {
	Email string `json:"email"`
}

type DisableUserResponse struct {
	Success bool `json:"success"`
}

// DisableUser set user status to active
func (h Handler) DisableUser() http.HandlerFunc {
	return rest.ErrResponseWrapper(func(w http.ResponseWriter, r *http.Request) error {
		req, err := request.BindJSON[DisableUserRequest](r.Body)
		if err != nil {
			return err
		}

		if err := h.userCtrl().DisableUser(r.Context(), req.Email); err != nil {
			return convertCtrlErr(err)
		}

		return respond.Ok(DisableUserResponse{
			Success: true,
		}).WriteJSON(w)
	})
}
