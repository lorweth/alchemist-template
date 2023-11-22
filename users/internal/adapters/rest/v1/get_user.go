package v1

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
	"github.com/virsavik/alchemist-template/users/internal/core/domain"
	"github.com/virsavik/alchemist-template/users/internal/core/ports"
	"github.com/virsavik/alchemist-template/users/internal/pkg/pagination"
)

func (hdl UserHandler) GetUser() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// Get query params
		req, err := httpio.URLQuery[getUserRequest](r, "query")
		if err != nil {
			return wrapBadRequestError(err, "invalid_request")
		}

		if err := req.IsValid(); err != nil {
			return err
		}

		// Get users
		users, total, err := hdl.svc.GetAll(ctx, ports.GetUserInput{
			Email: req.Email,
			CreatedAt: ports.Period{
				From: req.CreatedAt.From,
				To:   req.CreatedAt.To,
			},
			Pagination: req.Pagination,
		})
		if err != nil {
			return convertServiceError(err)
		}

		// Write response
		httpio.WriteJSON(w, r, httpio.Response[getUserResponse]{
			Status: http.StatusOK,
			Body: getUserResponse{
				Data: users,
				Meta: queryMeta{
					Total:       total,
					CurrentPage: req.Pagination.Page,
					Size:        req.Pagination.Size,
				},
			},
		})

		return nil
	})
}

type getUserRequest struct {
	Email      string           `json:"email"`
	CreatedAt  periodTime       `json:"created_at"`
	Pagination pagination.Input `json:"pagination"`
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

	if err := req.Pagination.IsValid(); err != nil {
		return wrapBadRequestError(err, "invalid_pagination")
	}

	return nil
}

type getUserResponse struct {
	Data []domain.User `json:"data"`
	Meta queryMeta     `json:"meta"`
}

type queryMeta struct {
	CurrentPage int   `json:"current_page"`
	Size        int   `json:"size"`
	Total       int64 `json:"total"`
}

type periodTime struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
