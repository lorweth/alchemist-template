package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/virsavik/alchemist-template/pkg/iam"
	"github.com/virsavik/alchemist-template/pkg/logger"
	"github.com/virsavik/alchemist-template/pkg/rest/respond"
)

// Authenticator is a middleware function that handles JWT authentication using Auth0's RS256 and JWKS flow.
// It extracts the JWT token from the Authorization header, validates it using the provided IAM validator,
// and sets the user profile information in the request context. If authentication fails, it responds with
// an unauthorized status and an error message.
//
// This middleware also enriches the request context with user-related information, such as the user ID,
// extracted from the authenticated token.
//
// For more information on the RS256 and JWKS flow, refer to: https://auth0.com/blog/navigating-rs256-and-jwks/#Verifying-a-JWT-using-the-JWKS-endpoint
func Authenticator(iamValidator iam.Validator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromCtx(ctx)

			p, err := getUserProfileFromRequest(r, iamValidator)
			if err != nil {
				log.Infof("user authenticate error: %v", err)
				respond.Unauthorized(respond.Message{Name: "invalid_token"}).WriteJSON(ctx, w)

				return
			}

			// Set user profile in the request context
			ctx = iam.SetInCtx(r.Context(), p)

			// Enrich context with user ID for logging purposes
			ctx = logger.SetInCtx(ctx, log.With(logger.String("enduser.id", p.ID)))

			// Pass control to the next handler with the enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// getUserProfileFromRequest extracts the JWT token from the request's authorization header,
// validates it, and returns the user profile associated with the token.
func getUserProfileFromRequest(r *http.Request, iamValidator iam.Validator) (iam.UserProfile, error) {
	// Extract the JWT from the request's authorization header.
	tokenRaw, err := getFromAuthHeader(r)
	if err != nil {
		return iam.UserProfile{}, err
	}

	// Let secure process the request. If it returns an error,
	// that indicates the request should not continue.
	parsedToken, err := iam.ParseJWT(iamValidator, tokenRaw)
	if err != nil {
		return iam.UserProfile{}, err
	}

	// Get UserProfile from token
	p, err := iam.GetUserProfile(parsedToken)
	if err != nil {
		return iam.UserProfile{}, err
	}
	return p, nil
}

// getFromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func getFromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("authorization header format must be bearer {token}")
	}

	return authHeaderParts[1], nil
}
