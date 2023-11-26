package jwt

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/iam"
)

type Provider struct {
	alg             jwa.SignatureAlgorithm
	sign            string // private key
	verifier        jwt.ParseOption
	validateOptions []jwt.ValidateOption
}

func NewProvider(cfg config.AppConfig, alg jwa.SignatureAlgorithm, sign string, options []jwt.ValidateOption) iam.Provider {
	return Provider{
		alg:             alg,
		sign:            sign,
		verifier:        jwt.WithKey(alg, sign),
		validateOptions: options,
	}
}

func (p Provider) Verifier() jwt.ParseOption {
	return p.verifier
}
