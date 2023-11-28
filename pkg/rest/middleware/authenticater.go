package middleware

import (
	"net/http"

	"github.com/virsavik/alchemist-template/pkg/iam"
	"github.com/virsavik/alchemist-template/pkg/iam/extractor"
	"github.com/virsavik/alchemist-template/pkg/iam/validator"
	"github.com/virsavik/alchemist-template/pkg/logger"
	"github.com/virsavik/alchemist-template/pkg/rest/httpio"
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
func Authenticator(validator validator.Validator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromCtx(ctx)

			p, err := getUserProfileFromRequest(r, validator)
			if err != nil {
				log.Infof("user authenticate error: %v", err)

				vErr := convertValidatorError(err)
				httpio.WriteJSON(w, r, httpio.Response[httpio.Message]{
					Status: vErr.Status,
					Body: httpio.Message{
						Code: vErr.Code,
						Desc: vErr.Desc,
					},
				})

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
func getUserProfileFromRequest(r *http.Request, validator validator.Validator) (iam.UserProfile, error) {
	// Extract the JWT from the request's authorization header.
	tokenRaw, err := extractor.AuthHeader(r)
	if err != nil {
		return iam.UserProfile{}, err
	}

	// Let secure process the request. If it returns an error,
	// that indicates the request should not continue.
	parsedToken, err := validator.ValidateToken(r.Context(), tokenRaw)
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

var (
	ErrUnauthorized = httpio.Error{Status: http.StatusUnauthorized, Code: "unauthorized", Desc: "token in unauthorized"}
	ErrExpired      = httpio.Error{Status: http.StatusUnauthorized, Code: "token_expired", Desc: "token is expired"}
	ErrNBFInvalid   = httpio.Error{Status: http.StatusUnauthorized, Code: "nbf_invalid", Desc: "token nbf validation failed"}
	ErrIATInvalid   = httpio.Error{Status: http.StatusUnauthorized, Code: "iat_invalid", Desc: "token iat validation failed"}
)

// convertValidatorError will normalize the error message from the underlining
func convertValidatorError(err error) httpio.Error {
	switch err.Error() {
	case validator.ErrExpired.Error():
		return ErrExpired
	case validator.ErrIATInvalid.Error():
		return ErrIATInvalid
	case validator.ErrNBFInvalid.Error():
		return ErrNBFInvalid
	default:
		return ErrUnauthorized
	}
}
