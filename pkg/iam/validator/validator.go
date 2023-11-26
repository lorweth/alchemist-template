package validator

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"

	"github.com/virsavik/alchemist-template/pkg/iam"
)

type validator struct {
	provider        iam.Provider
	validateOptions []jwt.ValidateOption
}

func New(provider iam.Provider, options ...jwt.ValidateOption) Validator {
	return validator{
		provider:        provider,
		validateOptions: options,
	}
}

func (v validator) ValidateToken(ctx context.Context, tokenRaw string) (jwt.Token, error) {
	// Parse token from raw
	token, err := jwt.Parse([]byte(tokenRaw), v.provider.Verifier(), jwt.WithValidate(false))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Return error unauthorized if token is nil
	if token == nil {
		return nil, ErrUnauthorized
	}

	// Validate token
	if err := jwt.Validate(token, v.validateOptions...); err != nil {
		return token, convertValidatorError(err)
	}

	return token, nil
}
