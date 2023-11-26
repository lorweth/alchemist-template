package iam

import (
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// GetUserProfile returns UserProfile by given token
func GetUserProfile(token jwt.Token) (UserProfile, error) {
	sub, ok := token.Get("sub")
	if !ok || sub == "" {
		return UserProfile{}, ErrTokenSubNotFound
	}

	id, ok := sub.(string)
	if !ok {
		return UserProfile{}, ErrTokenInvalid
	}

	return UserProfile{
		ID: id,
	}, nil
}
