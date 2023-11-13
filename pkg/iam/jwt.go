package iam

import (
	"fmt"

	"github.com/form3tech-oss/jwt-go"
)

// ParseJWT is a function that parses a JWT (JSON Web Token) using the provided validator.
// It verifies the token's audience, issuer, and signature based on the information from the validator.
//
// If the token is valid, it returns the parsed JWT token; otherwise, it returns an error indicating the reason for failure.
func ParseJWT(validator Validator, tokenRaw string) (*jwt.Token, error) {
	// If the tokenRaw is empty
	if tokenRaw == "" {
		return nil, errTokenIsBlank
	}

	// Decode the JWT and grab the kid property from the header.
	parsedToken, err := jwt.Parse(tokenRaw, func(token *jwt.Token) (interface{}, error) {
		// Ensure the JWT contains the expected audience, issuer, expiration, etc.
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errConvertToClaims
		}

		if !mapClaims.VerifyAudience(validator.GetAudience(), true) {
			return nil, errInvalidAudience
		}

		if !mapClaims.VerifyIssuer(validator.GetIssuer(), true) {
			return nil, errInvalidIssuer
		}

		// TODO: use alg HS256 for client JWT
		//if token.Header["alg"] != jwt.SigningMethodHS256.Alg() {
		//	return nil, fmt.Errorf("token algorithm invalid, expected: %s, actual: %s", jwt.SigningMethodHS256.Alg(), token.Header["alg"])
		//}

		kid, exists := token.Header["kid"]
		if !exists {
			return nil, errKidNotFound
		}

		// Find the signature verification key in the filtered JWKS with a matching kid property.
		// Using the x5c property build a certificate which will be used to verify the JWT signature
		// (already build when DownloadSigningKeysPolling)
		key, err := validator.GetSigningKey(kid.(string))
		if err != nil {
			return nil, fmt.Errorf("verify jwt error %w", err)
		}

		return key, nil

	})
	if err != nil {
		return nil, fmt.Errorf("parse token error %w", err)
	}

	// Check if the parsed tokenRaw is valid...
	if !parsedToken.Valid {
		return nil, errTokenInvalid
	}

	return parsedToken, nil
}
