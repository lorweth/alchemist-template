package validator

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Validator interface {
	ValidateToken(ctx context.Context, tokenRaw string) (jwt.Token, error)
}

type KeyFunc func(ctx context.Context) jwk.Set
