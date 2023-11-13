package iam

import (
	"errors"
	"fmt"
)

var (
	errSigningKeySetEmpty = errors.New("the jwks endpoint did not contain any signature verification keys")
	errPublicKeyNotFound  = errors.New("public key not found")
	errTokenSubNotFound   = errors.New("token sub missing")
	errConvertToClaims    = errors.New("cannot convert to map claims")
	errTokenIsBlank       = errors.New("jwt raw is blank")
	errInvalidAudience    = errors.New("invalid audience")
	errInvalidIssuer      = errors.New("invalid issuer")
	errKidNotFound        = errors.New("token header not contain kid")
	errTokenInvalid       = errors.New("token invalid")
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
