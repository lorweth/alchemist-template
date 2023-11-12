package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// JWKSClient refer from https://auth0.com/blog/navigating-rs256-and-jwks/#Verifying-a-JWT-using-the-JWKS-endpoint
type JWKSClient struct {
	jwksURI string
}

func newJWKSClient(jwksURI string) JWKSClient {
	return JWKSClient{
		jwksURI: jwksURI,
	}
}

func (jc JWKSClient) getJWKS(ctx context.Context) (jwks, error) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Get(jc.jwksURI)
	if err != nil {
		return jwks{}, fmt.Errorf("get jwks error: %w", err)
	}
	defer func() {
		span := trace.SpanFromContext(ctx)
		if err := resp.Body.Close(); err != nil {
			span.RecordError(err)
		}
	}()

	var keySets jwks
	if err := json.NewDecoder(resp.Body).Decode(&keySets); err != nil {
		return jwks{}, fmt.Errorf("decode jwks error %w", err)
	}

	return keySets, nil
}

func (jc JWKSClient) getSigningKeys(ctx context.Context) (map[string]string, error) {
	keySet, err := jc.getJWKS(ctx)
	if err != nil {
		return nil, err
	}

	if len(keySet.Keys) == 0 {
		return nil, errors.New("the jwks endpoint did not contain any keys")
	}

	signingKeys := make(map[string]string)
	for _, key := range keySet.Keys {
		// JWK property `use` determines the JWK is for signature verification
		// We are only supporting RSA (RS256)
		// The `kid` must be present to be useful for later
		// x5c or (n and e) has useful public keys
		if key.Use == "sig" && key.Kty == "RSA" && key.Kid != "" &&
			len(key.X5c) > 0 || (key.N != "" && key.E != "") {
			signingKeys[key.Kid] = fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", key.X5c[0])
		}
	}

	// If at least one signing key doesn't exist we have a problem... Kaboom.
	if len(signingKeys) == 0 {
		return nil, errors.New("the jwks endpoint did not contain any signature verification keys")
	}

	return signingKeys, nil
}
