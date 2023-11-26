package jwks

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/logger"
)

const (
	scopeName              = "github.com/virsavik/alchemist-template/pkg/iam/jwks/provider.go"
	defaultRefreshInterval = 15 * time.Minute
)

type CacheProvider struct {
	logger          logger.Logger
	tracer          trace.Tracer
	client          *http.Client
	jwkURI          string
	mu              sync.RWMutex
	refreshInterval time.Duration
	cache           jwk.Set // Default is nil
}

func NewProvider(cfg config.AppConfig, opts ...ProviderOption) (*CacheProvider, error) {
	log, err := logger.New(logger.Config{
		Environment: cfg.Environment,
	})
	if err != nil {
		return nil, err
	}

	p := &CacheProvider{
		logger:          log,
		tracer:          otel.Tracer(scopeName),
		jwkURI:          getJwkURI(cfg.IAM.Tenant),
		client:          &http.Client{},
		refreshInterval: defaultRefreshInterval,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func (cp *CacheProvider) Verifier() jwt.ParseOption {
	return jwt.WithKeySet(cp.GetKeySet())
}

func (cp *CacheProvider) GetKeySet() jwk.Set {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	return cp.cache
}

func (cp *CacheProvider) fetchJWKS(ctx context.Context) (jwk.Set, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Create a new GET request with the context
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cp.jwkURI, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request error: %w", err)
	}

	// Call request
	resp, err := cp.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("making request error: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			cp.logger.Errorf(err, "close response body")
		}
	}()

	// Parse to jwk set
	set, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse jwks error: %w", err)
	}

	return set, nil
}

func (cp *CacheProvider) cacheJWKS(jwks jwk.Set) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	cp.cache = jwks
}

func (cp *CacheProvider) RefreshLoop(ctx context.Context) error {
	log := cp.logger
	ticker := time.NewTicker(defaultRefreshInterval)

	go func() {
		for {
			// Start tracer
			ctx, span := cp.tracer.Start(
				context.Background(),
				"fetchJWKS",
				trace.WithAttributes(
					semconv.HTTPRouteKey.String(cp.jwkURI),
				),
				trace.WithSpanKind(trace.SpanKindClient),
			)

			// Fetch JWKS
			set, err := cp.fetchJWKS(ctx)
			if err != nil {
				log.Errorf(err, "fetch jwks error")
				span.RecordError(err, trace.WithStackTrace(true))
			} else {
				log.Infof("fetch jwks complete successfully")
				span.AddEvent("fetch successfully")
			}

			// Cache value
			cp.cacheJWKS(set)

			// End span
			// Must be before sleep process
			span.End()

			// Sleep in d second
			<-ticker.C
		}
	}()

	// Wait for stop process
	<-ctx.Done()

	log.Infof("Stop download signing key")

	// Stop ticker
	ticker.Stop()

	return nil
}

// getJwkURI return auth0 jwks_uri from TENANT
func getJwkURI(tenant string) string {
	return fmt.Sprintf("https://%s.auth0.com/.well-known/jwks.json", tenant)
}
