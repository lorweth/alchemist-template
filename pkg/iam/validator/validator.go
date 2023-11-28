package validator

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
)

type validator struct {
	validateOptions []jwt.ValidateOption
	parser          Parser
}

func New(parser Parser, options ...jwt.ValidateOption) Validator {
	return validator{
		parser:          parser,
		validateOptions: options,
	}
}

func (v validator) ValidateToken(ctx context.Context, tokenRaw string) (jwt.Token, error) {
	// Parse token from raw
	token, err := v.parser.Parse(tokenRaw)
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
