package users

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/virsavik/alchemist-template/pkg/postgres"
	"github.com/virsavik/alchemist-template/pkg/rest/middleware"
	"github.com/virsavik/alchemist-template/pkg/system"
	"github.com/virsavik/alchemist-template/users/internal/adapters/repository"
	"github.com/virsavik/alchemist-template/users/internal/adapters/repository/generator"
	v1 "github.com/virsavik/alchemist-template/users/internal/adapters/rest/v1"
	"github.com/virsavik/alchemist-template/users/internal/core/services"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	// Init sonyflake id generator
	generator.InitIDGenerator()

	store := repository.New(postgres.Trace(svc.DB()))
	userService := services.NewUserService(store)
	userHandler := v1.NewUserHandler(userService)

	setupRoutes(svc, *userHandler)

	return nil
}

func setupRoutes(svc system.Service, hdl v1.UserHandler) {
	svc.Mux().Use(middleware.Logger(svc.Logger()))
	svc.Mux().Use(middleware.Recover())
	svc.Mux().Use(middleware.OtelTracer())

	svc.Mux().Route("/users", func(v1 chi.Router) {
		//v1.Use(middleware.Authenticator(svc.IAMValidator()))

		v1.Post("/", hdl.CreateUser())
		v1.Delete("/{id}", hdl.DeleteUser())
	})
}
