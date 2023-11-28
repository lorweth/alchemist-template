package jwks

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
	"github.com/virsavik/alchemist-template/pkg/iam/oidc"
	"github.com/virsavik/alchemist-template/pkg/logger"
)

const (
	defaultRefreshInterval = 15 * time.Minute
	defaultRequestTimeout  = 10 * time.Second
)

type CacheProvider struct {
	issuer          url.URL
	jwksURI         url.URL
	audience        string
	refreshInterval time.Duration
	client          http.Client
	mu              sync.RWMutex
	logger          logger.Logger
	cache           jwk.Set // Init when first time fetch value
}

func NewProvider(tenant string, audience string, opts ...Option) (*CacheProvider, error) {
	issuer, err := oidc.GetIssuerFromTenant(tenant)
	if err != nil {
		return nil, err
	}

	p := &CacheProvider{
		issuer:          issuer,
		jwksURI:         oidc.GetJWKSURI(issuer),
		audience:        audience,
		refreshInterval: defaultRefreshInterval,
		client:          http.Client{Timeout: defaultRequestTimeout},
		logger:          logger.NewNoop(), // Default logger
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}
func (c *CacheProvider) Parse(tokenRaw string) (jwt.Token, error) {
	return jwt.Parse(
		[]byte(tokenRaw),
		jwt.WithKeySet(c.GetPublicKeys()),
		jwt.WithValidate(false), // validate by jwt.Validate
		jwt.WithAudience(c.audience),
	)
}

func (c *CacheProvider) GetPublicKeys() jwk.Set {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache
}

func (c *CacheProvider) fetchPublicKeys(ctx context.Context) (jwk.Set, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURI.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		resp.Body.Close()
	}()

	jwks, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return jwks, nil
}

func (c *CacheProvider) storePublicKeys(jwks jwk.Set) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = jwks
}

func (c *CacheProvider) FetchLoop(ctx context.Context) error {
	log := c.logger
	log.Infof("Starting fetch JWKS loop...")

	ticker := time.NewTicker(5 * time.Minute)

	// Intentionally not having a separate func to control polling as we don't need to care about graceful shutdown for this poller or explicitly handle cleanup
	go func() {
		for {
			log.Infof("Attempting download and caching JWKS...")

			newCtx := logger.NewCtx(ctx)

			// Fetch JWKS
			jwks, err := c.fetchPublicKeys(newCtx)
			if err != nil {
				log.Errorf(err, "download and caching failed")
			} else {
				log.Infof("Download and caching complete")
			}

			// Cache JWKS
			c.storePublicKeys(jwks)

			<-ticker.C
		}
	}()

	<-ctx.Done()

	log.Infof("Stopping JWKS fetch loop...")

	ticker.Stop()
	return nil
}
