package system

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/logger"
	"github.com/virsavik/alchemist-template/pkg/waiter"
)

// Service representing an application service
type Service interface {
	Config() config.AppConfig
	DB() *sql.DB
	Mux() *chi.Mux
	Logger() logger.Logger
	Waiter() waiter.Waiter
}

// Module representing an application module
type Module interface {
	Startup(context.Context, Service) error
}
