package middlewares

import (
	rptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rs/zerolog"
	"net/http"
)

func ExchangeSubjectForUser(logger zerolog.Logger, logic provider.Controller) func(next http.Handler) http.Handler {
	logger = logger.With().Str("module", "middleware").Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			subject, ok := rptrCtx.Subject(r.Context())
			if !ok {
				logger.Error().Msgf("subject not on context")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userName, doesExist, err := logic.UserExists(r.Context(), subject)
			if err != nil {
				logger.Error().Err(err).Msgf("while checking if user exists for subject: %s", subject)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !doesExist {
				logger.Warn().Msgf("no user for subject: %s", subject)
				w.WriteHeader(http.StatusNotFound)
				// TODO: Not found handler
				return
			}
			logger.Debug().Msgf("username from repo: %s", userName)
			r = r.WithContext(rptrCtx.AddUsername(r.Context(), userName))
			next.ServeHTTP(w, r)
		})
	}
}
