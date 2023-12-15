package oidc

import (
	"fmt"
	"net/url"
)

// GetIssuerFromTenant gets the issuer url from tenant.
func GetIssuerFromTenant(tenant string) (url.URL, error) {
	issuer, err := url.Parse(fmt.Sprintf("https://%s.auth0.com", tenant))
	if err != nil {
		return url.URL{}, fmt.Errorf("could not parse issuer from tenant: %w", err)
	}

	return *issuer, nil
}

// GetJwksURI gets the jwks uri from issuer.
func GetJwksURI(issuer url.URL) url.URL {
	return *issuer.JoinPath(".well-known/jwks.json")
}
