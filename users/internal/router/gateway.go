package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/virsavik/alchemist-template/pkg/di"
	"github.com/virsavik/alchemist-template/users/internal/handler/rest"
)

func RegisterGateway(ctn di.Container, mux *chi.Mux) {
	const apiRoot = "/api/users"

	hdl := rest.New(ctn)

	mux.Route(apiRoot, func(r chi.Router) {
		r.Post("/", hdl.RegisterUser())
		r.Put("/enable", hdl.EnableUser())
		r.Put("/disable", hdl.DisableUser())
	})
}
