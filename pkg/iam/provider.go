package iam

import (
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Provider interface {
	Verifier() jwt.ParseOption
}
