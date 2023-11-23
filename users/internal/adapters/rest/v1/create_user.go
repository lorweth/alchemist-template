package v1

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
)

func (hdl UserHandler) CreateUser() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// Decode request
		req, err := httpio.BindJSON[createUserRequest](r.Body)
		if err != nil {
			return err
		}

		// Validate request
		if err := req.IsValid(); err != nil {
			return err
		}

		// Create user
		user, err := hdl.svc.Create(ctx, domain.User{
			Email: req.Email,
		})
		if err != nil {
			return convertServiceError(err)
		}

		// Write response
		httpio.WriteJSON(w, r, httpio.Response[domain.User]{
			Status: http.StatusOK,
			Body:   user,
		})

		return nil
	})
}

type createUserRequest struct {
	Email string `json:"email"`
}

func (req *createUserRequest) IsValid() error {
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return errEmailCannotBeBlank
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return wrapBadRequestError(err, "invalid_email")
	}

	return nil
}
