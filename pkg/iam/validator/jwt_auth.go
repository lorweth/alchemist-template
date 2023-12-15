package validator

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	HS256 = jwa.HS256
)

type JWTAuth struct {
	issuer   string
	audience string
	alg      jwa.SignatureAlgorithm
	signKey  []byte
}

func NewJWTAuth(alg jwa.SignatureAlgorithm, signKey []byte, issuer string, audience string) JWTAuth {
	return JWTAuth{
		alg:      alg,
		signKey:  signKey,
		issuer:   issuer,
		audience: audience,
	}
}

func (ja JWTAuth) Parse(tokenRaw string) (jwt.Token, error) {
	return ja.Decode(tokenRaw)
}

func (ja JWTAuth) Encode(claims Claims) (jwt.Token, string, error) {
	// Create new token from claims
	tk, err := claims.TokenBuilder().Build()
	if err != nil {
		return nil, "", fmt.Errorf("create token err: %w", err)
	}

	// Create a signed JWT token serialized
	tokenStr, err := jwt.Sign(tk, jwt.WithKey(ja.alg, ja.signKey))
	if err != nil {
		return nil, "", fmt.Errorf("create signed jwt token error: %w", err)
	}

	return tk, string(tokenStr), nil
}

func (ja JWTAuth) Decode(tokenRaw string) (jwt.Token, error) {
	return jwt.Parse([]byte(tokenRaw), jwt.WithKey(ja.alg, ja.signKey))
}
