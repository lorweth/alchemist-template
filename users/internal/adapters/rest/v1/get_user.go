package v1

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
)

func (hdl UserHandler) GetUser() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// Get query params
		req, err := httpio.URLQuery[getUserRequest](r, "q")
		if err != nil {
			return wrapBadRequestError(err, "invalid_request")
		}

		if err := req.IsValid(); err != nil {
			return err
		}

		// Delete user
		users, err := hdl.svc.GetAll(ctx, ports.GetUserInput{
			Email: req.Email,
			CreatedAt: ports.Period{
				From: req.CreatedAt.From,
				To:   req.CreatedAt.To,
			},
		})
		if err != nil {
			return convertServiceError(err)
		}

		// Write response
		httpio.WriteJSON(w, r, httpio.Response[getUserResponse]{
			Status: http.StatusOK,
			Body: getUserResponse{
				Data: users,
			},
		})

		return nil
	})
}

type getUserRequest struct {
	Email     string     `json:"email"`
	CreatedAt periodTime `json:"created_at"`
}

func (req getUserRequest) IsValid() error {
	req.Email = strings.TrimSpace(req.Email)
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return wrapBadRequestError(err, "invalid_email")
		}
	}

	if req.CreatedAt.To.Before(req.CreatedAt.From) {
		return errToTimeCannotBeforeFromTime
	}

	return nil
}

type getUserResponse struct {
	Data []domain.User `json:"data"`
}

type periodTime struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
