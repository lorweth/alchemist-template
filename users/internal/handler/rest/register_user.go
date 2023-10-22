package rest

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/rest"
	"github.com/virsavik/alchemist-template/pkg/rest/request"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
	"github.com/virsavik/alchemist-template/users/internal/model"
)

type RegisterUserRequest struct {
	Email string `json:"email"`
}

type RegisterUserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// RegisterUser register user api
func (h Handler) RegisterUser() http.HandlerFunc {
	return rest.ErrResponseWrapper(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := request.BindJSON[RegisterUserRequest](r.Body)
		if err != nil {
			return err
		}

		user := model.User{
			Email: req.Email,
		}
		if err := h.userCtrl().RegisterUser(ctx, &user); err != nil {
			return convertCtrlErr(err)
		}

		return respond.Ok(RegisterUserResponse{
			ID:    user.ID,
			Email: user.Email,
		}).WriteJSON(w)
	})
}
