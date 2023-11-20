package iam

import (
	"context"
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/logger"
)

const (
	defaultDownloadJWKSDuration = 15 * time.Minute
)

// authZeroValidator is a struct implementing the Validator interface, specifically designed
// for authenticating JWTs using the Auth0 service.
type authZeroValidator struct {
	issuer           string
	audience         string
	jwksClient       *jwksClient
	tracer           trace.Tracer
	logger           logger.Logger
	cachedSigningKey map[string]*rsa.PublicKey
	mu               *sync.RWMutex
}

// New creates a new instance of authZeroValidator using the provided AppConfig configuration.
// It initializes the necessary components such as logger, issuer, audience, JWKS client, tracer, and mutex.
func New(cfg config.AppConfig) (Validator, error) {
	log, err := logger.New(logger.Config{
		Environment: cfg.Environment,
	})
	if err != nil {
		return nil, err
	}

	issuer := fmt.Sprintf("https://%s.auth0.com/", cfg.IAM.Tenant)

	return &authZeroValidator{
		issuer:           issuer,
		audience:         cfg.IAM.Audience,
		jwksClient:       newJWKSClient(fmt.Sprintf("%s.well-known/jwks.json", issuer), defaultDownloadJWKSTimeout),
		tracer:           otel.Tracer("auth0Validator"),
		logger:           log,
		cachedSigningKey: nil, // will be initialized after first call
		mu:               &sync.RWMutex{},
	}, nil
}

// GetIssuer returns the issuer URL configured for the Auth0 Validator.
func (v *authZeroValidator) GetIssuer() string {
	return v.issuer
}

// GetAudience returns the expected audience configured for the Auth0 Validator.
func (v *authZeroValidator) GetAudience() string {
	return v.audience
}

// downloadSigningKeys fetches and caches the JWKS (JSON Web Key Set) from the configured endpoint.
func (v *authZeroValidator) downloadSigningKeys(ctx context.Context) error {
	// Get JWKs
	keySet, err := v.jwksClient.getJWKs(ctx)
	if err != nil {
		return wrapError(err, "download signing keys")
	}

	// Parse JWKs to rsa.PublicKey
	publicKeys, err := v.jwksClient.parsePublicKeys(ctx, keySet)
	if err != nil {
		return err
	}

	// Cached public keys
	v.mu.Lock()
	defer v.mu.Unlock()
	// Cached signing keys
	v.cachedSigningKey = publicKeys

	return nil
}

// GetSigningKey retrieves a cached signing key based on the provided key ID (kid).
// If the key is not found, it returns err ErrPublicKeyNotFound.
func (v *authZeroValidator) GetSigningKey(kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	key, exists := v.cachedSigningKey[kid]
	if !exists {
		return nil, ErrPublicKeyNotFound
	}

	return key, nil
}

// DownloadSigningKeysPolling continuously polls the JWKS endpoint at regular intervals.
// It runs as a goroutine and stops when the provided context is canceled.
// It logs success and error events, and the polling frequency is controlled by defaultDownloadJWKSDuration.
func (v *authZeroValidator) DownloadSigningKeysPolling(ctx context.Context) error {
	log := v.logger
	d := defaultDownloadJWKSDuration
	ticker := time.NewTicker(d)

	go func() {
		for {
			ctx, span := v.tracer.Start(
				context.Background(),
				"downloadSigningKeys",
				trace.WithAttributes(
					semconv.HTTPRouteKey.String(v.jwksClient.uri),
				),
				trace.WithSpanKind(trace.SpanKindClient),
			)

			if err := v.downloadSigningKeys(ctx); err != nil {
				log.Errorf(err, "download singing keys error")
				span.AddEvent("downloadSigningKeys error", trace.WithAttributes(
					attribute.String("error", err.Error()),
				))
			} else {
				log.Infof("download signing keys complete")
				span.AddEvent("downloadSigningKeys finished")
			}

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
