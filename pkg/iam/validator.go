package iam

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/logger"
)

const (
	defaultDownloadJWKSDuration = 5 * time.Second
)

type authZeroValidator struct {
	domain           string
	audience         string
	jwksURI          string
	tracer           trace.Tracer
	logger           logger.Logger
	cachedSigningKey map[string]string
	mu               *sync.RWMutex
}

func New(cfg config.AppConfig) (Validator, error) {
	log, err := logger.New(logger.LogConfig{
		Environment: cfg.Environment,
	})
	if err != nil {
		return nil, err
	}

	return &authZeroValidator{
		domain:           cfg.IAM.Domain,
		audience:         cfg.IAM.Audience,
		jwksURI:          fmt.Sprintf("%s/.well-known/jwks.json", strings.Trim(cfg.IAM.Domain, "/")),
		tracer:           otel.Tracer("authZeroValidator"),
		logger:           log,
		cachedSigningKey: nil, // will be initialized after first call
		mu:               &sync.RWMutex{},
	}, nil
}

func (v *authZeroValidator) GetDomain() string {
	return v.domain
}

func (v *authZeroValidator) GetAudience() string {
	return v.audience
}

func (v *authZeroValidator) downloadSigningKeys(ctx context.Context) error {
	client := newJWKSClient(v.jwksURI)

	keys, err := client.getSigningKeys(ctx)
	if err != nil {
		return err
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	// Cached signing keys
	v.cachedSigningKey = keys

	return nil
}

func (v *authZeroValidator) getSigningKey(kid string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	key, exists := v.cachedSigningKey[kid]
	if !exists {
		return "", errors.New("singing key not found")
	}

	return key, nil
}

func (v *authZeroValidator) VerifyJWT(kid string) (*rsa.PublicKey, error) {
	signingKey, err := v.getSigningKey(kid)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(signingKey))
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (v *authZeroValidator) PollingDownloadSigningKeys(ctx context.Context, d time.Duration) {
	if d == 0 {
		d = defaultDownloadJWKSDuration
	}

	timer := time.NewTimer(d)
	log := v.logger

	go func() {
		for {
			log.Info("auth0 signing keys downloading")
			ctx, span := v.tracer.Start(
				context.Background(),
				"downloadSigningKeys",
				trace.WithAttributes(
					semconv.HTTPRouteKey.String(v.jwksURI),
				),
				trace.WithSpanKind(trace.SpanKindClient),
			)

			if err := v.downloadSigningKeys(ctx); err != nil {
				log.Error(err, "download singing keys error")
				span.AddEvent("downloadSigningKeys error", trace.WithAttributes(
					attribute.String("error", err.Error()),
				))
			}

			log.Info("Download signing keys finished")
			span.AddEvent("downloadSigningKeys finished")

			span.End()

			// Sleep in d second
			<-timer.C

			// Reset time with d section
			timer.Reset(d)
		}
	}()
}
