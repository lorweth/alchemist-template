package jwks

import "net/http"

// ProviderOption is how options for the CacheProvider are set up.
type ProviderOption func(*CacheProvider)

// WithCustomClient will set a custom *http.Client on the *CacheProvider
func WithCustomClient(c *http.Client) ProviderOption {
	return func(p *CacheProvider) {
		p.client = c
	}
}
