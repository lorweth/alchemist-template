package middleware

import (
	"errors"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"

	"github.com/virsavik/alchemist-template/pkg/iam"
	"github.com/virsavik/alchemist-template/pkg/logger"
)

func Auth(log logger.Logger, iamValidator iam.Validator) func(next http.Handler) http.Handler {
	mw := newJWTMiddleware(log, iamValidator)

	return func(next http.Handler) http.Handler {
		return mw.Handler(next)
	}
}

type jwtMiddleware struct {
	logger logger.Logger
	*jwtmiddleware.JWTMiddleware
}

func newJWTMiddleware(log logger.Logger, iamValidator iam.Validator) jwtMiddleware {
	jwtMW := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := iamValidator.GetAudience()
			if !token.Claims.(jwt.MapClaims).VerifyAudience(aud, false) {
				return nil, errors.New("invalid audience")
			}

			// Verify 'iss' claim
			iss := iamValidator.GetDomain()
			if !token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false) {
				return nil, errors.New("invalid issuer")
			}

			kid, exists := token.Header["kid"]
			if !exists {
				return nil, errors.New("token header not contain kid")
			}

			publicKey, err := iamValidator.VerifyJWT(kid.(string))
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return jwtMiddleware{
		logger:        log,
		JWTMiddleware: jwtMW,
	}
}

func (m jwtMiddleware) CheckJWT(w http.ResponseWriter, r *http.Request) error {
	if err := m.JWTMiddleware.CheckJWT(w, r); err != nil {
		m.logger.Error(err, "check jwt error")
	}

	m.logger.Info("jwt valid")

	return nil
}
