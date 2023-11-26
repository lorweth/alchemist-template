package validator

import (
	"errors"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrUnauthorized = errors.New("token is unauthorized")
	ErrExpired      = errors.New("token is expired")
	ErrNBFInvalid   = errors.New("token nbf validation failed")
	ErrIATInvalid   = errors.New("token iat validation failed")
)

// convertValidatorError will normalize the error message from the underlining
// jwt library
func convertValidatorError(err error) error {
	switch err.Error() {
	case jwt.ErrTokenExpired().Error():
		return ErrExpired
	case jwt.ErrInvalidIssuedAt().Error():
		return ErrIATInvalid
	case jwt.ErrTokenNotYetValid().Error():
		return ErrNBFInvalid
	default:
		return ErrUnauthorized
	}
}
