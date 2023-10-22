package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"golang.org/x/sync/errgroup"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/logger"
	"github.com/virsavik/alchemist-template/pkg/waiter"
)

// System represents the core components that acts as a Facade in the application, including its configuration,
// database connection, HTTP request router, and a waiter for graceful shutdown. It serves
// as the central struct that encapsulates these essential elements for managing and
// controlling the application's behavior.
type System struct {
	cfg    config.AppConfig
	db     *sql.DB
	mux    *chi.Mux
	logger logger.Logger
	waiter waiter.Waiter
	tp     *sdktrace.TracerProvider
}

func New(cfg config.AppConfig) (*System, error) {
	s := &System{cfg: cfg}

	s.initWaiter()

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initOpenTelemetry(); err != nil {
		return nil, err
	}

	s.initLogger()

	s.initMux()

	return s, nil
}

func (s *System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) initDB() (err error) {
	s.db, err = sql.Open("pgx", s.cfg.PG.URI)
	return err
}

func (s *System) DB() *sql.DB {
	return s.db
}

func (s *System) initLogger() {
	var err error
	s.logger, err = logger.New(logger.LogConfig{
		Environment: s.cfg.Environment,
	})
	if err != nil {
		panic("init logger error")
	}
}

func (s *System) Logger() logger.Logger {
	return s.logger
}

func (s *System) initMux() {
	s.mux = chi.NewMux()
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

// initOpenTelemetry Initializes an OTLP exporter, and configures the corresponding trace
// and metric providers.
// copy from sample https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go
func (s *System) initOpenTelemetry() error {
	// Skip when config is not provided
	if s.cfg.Otel.ServiceName == "" || s.cfg.Otel.ExporterEndpoint == "" {
		return nil
	}

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(s.cfg.Otel.ServiceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up a trace exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(s.cfg.Otel.ExporterEndpoint),
	)
	if err != nil {
		return err
	}

	// Register the trace exporter with a TracerProvider, using
	// a batch span processor to aggregate spans before export.
	s.tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(s.tp)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	s.waiter.Cleanup(func() {
		if err := s.tp.Shutdown(context.Background()); err != nil {
			s.logger.Error(err, "ran into an issue shutting down the tracer provider")
		}
	})

	return nil
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Port)
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}
