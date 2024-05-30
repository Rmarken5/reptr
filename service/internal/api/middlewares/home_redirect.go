package middlewares

import (
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rs/zerolog"
	"net/http"
)

func HomeRedirect(logger zerolog.Logger, authenticator auth.Authentication) func(next http.Handler) http.Handler {
	logger = logger.With().Str("middleware", "HomeRedirect").Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Msg("home redirect middleware")
			logger.Info().Msg(r.URL.Path)

			authHeader := r.Header.Get("Authorization")
			token, err := parseBearerToken(authHeader)
			if err != nil {
				logger.Error().Err(err).Msg("while parsing auth header")
				next.ServeHTTP(w, r)
				return
			}

			idToken, err := authenticator.VerifyIDToken(r.Context(), token)
			if err != nil {
				logger.Error().Err(err).Msg("while verifying token")
				next.ServeHTTP(w, r)
				return
			}
			if idToken != nil {
				logger.Info().Msg("redirecting")
				http.Redirect(w, r, "/page/home", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
