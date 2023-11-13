package iam

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"go.opentelemetry.io/otel/trace"
)

var (
	defaultDownloadJWKSTimeout = 5 * time.Second
)

// jwksClient represents a client for interacting with a JWKS (JSON Web Key Set) endpoint.
type jwksClient struct {
	uri        string
	httpClient http.Client
}

// newJWKSClient creates a new JWKS client with the specified URI and timeout duration.
// If the timeout is not provided (set to zero), it defaults to defaultDownloadJWKSTimeout.
func newJWKSClient(jwksURI string, timeout time.Duration) *jwksClient {
	if timeout == 0 {
		timeout = defaultDownloadJWKSTimeout
	}

	return &jwksClient{
		uri: jwksURI,
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

// getJWKs retrieves the JWKS (JSON Web Key Set) from the configured URI.
func (jc jwksClient) getJWKs(ctx context.Context) (jwks, error) {
	resp, err := jc.httpClient.Get(jc.uri)
	if err != nil {
		return jwks{}, wrapError(err, "get jwks error")
	}
	defer func() {
		span := trace.SpanFromContext(ctx)
		if err := resp.Body.Close(); err != nil {
			span.RecordError(err)
		}
	}()

	var keySets jwks
	if err := json.NewDecoder(resp.Body).Decode(&keySets); err != nil {
		return jwks{}, wrapError(err, "decode jwks error")
	}

	return keySets, nil
}

// parsePublicKeys parses public keys from the given JWKS (JSON Web Key Set).
// It follows the RS256 (RSA Signature with SHA-256) algorithm for JWT verification.
//
// Reference: https://auth0.com/blog/navigating-rs256-and-jwks/#Verifying-a-JWT-using-the-JWKS-endpoint
func (jc jwksClient) parsePublicKeys(ctx context.Context, keySet jwks) (map[string]*rsa.PublicKey, error) {
	publicKeys := make(map[string]*rsa.PublicKey)
	for _, key := range keySet.Keys {
		// JWK property `use` determines the JWK is for signature verification
		// skip it if it is not "sig"
		if key.Use != "sig" {
			continue
		}

		// We are only supporting RSA (RS256)
		if key.Kty != "RSA" {
			continue
		}

		// The `kid` must be present to be useful for later
		if key.Kid == "" {
			continue
		}

		// x5c or (n and e) has useful public keys
		// Skip it if x5c is empty, and both n and e are blank
		if len(key.X5c) == 0 && (key.N == "" || key.E == "") {
			continue
		}

		// Parse to RSA public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", key.X5c[0])))
		if err != nil {
			return nil, err
		}

		// Add to map
		publicKeys[key.Kid] = publicKey
	}

	// If at least one signing key doesn't exist we have a problem... Kaboom.
	if len(publicKeys) == 0 {
		return nil, errSigningKeySetEmpty
	}

	return publicKeys, nil
}
