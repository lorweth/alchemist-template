package users

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"

	"github.com/virsavik/alchemist-template/pkg/di"
	"github.com/virsavik/alchemist-template/pkg/postgresotel"
	"github.com/virsavik/alchemist-template/pkg/rest/middleware"
	"github.com/virsavik/alchemist-template/pkg/system"
	"github.com/virsavik/alchemist-template/users/internal/constants"
	"github.com/virsavik/alchemist-template/users/internal/controller"
	"github.com/virsavik/alchemist-template/users/internal/repository"
	"github.com/virsavik/alchemist-template/users/internal/router"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	// Init dependency injection container
	container := di.New()

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})
	container.AddScoped(constants.UsersRepoKey, func(c di.Container) (any, error) {
		return repository.New(
			postgresotel.Trace(svc.DB()),
		), nil
	})

	userRegistered := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.UsersRegisteredCount,
	})
	container.AddScoped(constants.UsersCtrlKey, func(c di.Container) (any, error) {
		return controller.New(
			controller.NewUserController(c.Get(constants.UsersRepoKey).(repository.UserRepository)),
			userRegistered,
		), nil
	})

	// setup Driver adapters
	setupChiMiddleware(svc)

	setupMetricRoute(svc.Mux())

	router.RegisterGateway(container, svc.Mux())

	return nil
}

func setupChiMiddleware(svc system.Service) {
	svc.Mux().Use(middleware.RequestID())
	svc.Mux().Use(middleware.Logger(svc.Logger()))
	svc.Mux().Use(middleware.Otel(otel.Tracer("handler")))
	svc.Mux().Use(middleware.Auth(svc.Logger(), svc.IAMValidator()))
}

func setupMetricRoute(mux *chi.Mux) {
	mux.Method("GET", "/metrics", promhttp.Handler())
}
