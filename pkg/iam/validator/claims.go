package validator

import (
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

// RegisteredClaims represent standard claims
type RegisteredClaims struct {
	Issuer         string
	Subject        string
	Audience       string
	ExpirationTime time.Time
	NotBefore      time.Time
	IssuedAt       time.Time
	JwtID          string
}

func (clm RegisteredClaims) TokenBuilder() *jwt.Builder {
	builder := jwt.NewBuilder()

	builder.Subject(clm.Subject)
	builder.Expiration(clm.ExpirationTime.UTC())
	builder.Issuer(clm.Issuer)
	builder.Audience([]string{clm.Audience})

	if !clm.IssuedAt.IsZero() {
		builder.IssuedAt(clm.IssuedAt.UTC())
	}

	if !clm.NotBefore.IsZero() {
		builder.NotBefore(clm.NotBefore.UTC())
	}

	if clm.JwtID == "" {
		builder.JwtID(clm.JwtID)
	}

	return builder
}
