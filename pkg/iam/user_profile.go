package iam

import (
	"github.com/form3tech-oss/jwt-go"
)

// GetUserProfile returns UserProfile by given token
func GetUserProfile(token *jwt.Token) (UserProfile, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserProfile{}, errConvertToClaims
	}

	id, ok := mapClaims["sub"].(string)
	if !ok || id == "" {
		return UserProfile{}, errTokenSubNotFound
	}

	return UserProfile{
		ID: id,
	}, nil
}
