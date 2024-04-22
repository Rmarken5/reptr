package middlewares

import (
	"errors"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

func Authenticate(logger zerolog.Logger, authenticator auth.Authentication) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.With().Str("middleware", "Authenticate").Logger()
			logger.Info().Msg("authenticating")

			authHeader := r.Header.Get("Authorization")
			token, err := parseBearerToken(authHeader)
			if err != nil {
				logger.Error().Err(err).Msgf("Bad request - Invalid auth token")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			idToken, err := authenticator.VerifyIDToken(r.Context(), token)
			if err != nil {
				logger.Error().Err(err).Msgf("error authenticating %s: %v", token, err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			*r = *r.WithContext(reptrCtx.AddSubject(r.Context(), strings.TrimPrefix(idToken.Subject, "auth0|")))

			next.ServeHTTP(w, r)
		})
	}
}

func parseBearerToken(bearerToken string) (string, error) {
	authParts := strings.SplitN(bearerToken, " ", 2)

	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return "", errors.New("invalid bearer auth header")
	}
	return authParts[1], nil
}
