package iam

import (
	"errors"
	"fmt"
)

var (
	ErrSigningKeySetEmpty = errors.New("the jwks endpoint did not contain any signature verification keys")
	ErrPublicKeyNotFound  = errors.New("public key not found")
	ErrTokenSubNotFound   = errors.New("token sub missing")
	ErrConvertToClaims    = errors.New("cannot convert to map claims")
	ErrTokenIsBlank       = errors.New("jwt raw is blank")
	ErrInvalidAudience    = errors.New("invalid audience")
	ErrInvalidIssuer      = errors.New("invalid issuer")
	ErrKidNotFound        = errors.New("token header not contain kid")
	ErrTokenInvalid       = errors.New("token invalid")
)

// wrapError wraps error by given error and description
//
// Example usage:
//
//	err := somethingMethod()
//	wrappedErr := wrapError(err, "do something error")
func wrapError(err error, desc string) error {
	return fmt.Errorf(desc+": %w", err)
}
