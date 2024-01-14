package middlewares

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/service/internal/api"
	"github.com/rs/zerolog"
	"net/http"
)

func Session(logger zerolog.Logger, store *sessions.CookieStore) func(next http.Handler) http.Handler {
	logger = logger.With().Str("middleware", "Session").Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug().Msgf("session middleware invoked")

			session, _ := store.Get(r, api.CookieSessionID)
			logger.Debug().Msgf("Is session new? : %t", session.IsNew)
			if !session.IsNew {
				authToken, err := AuthFromSession(session.Values)
				if err != nil {
					logger.Error().Err(err).Msg("while getting auth token from session")
					http.Error(w, "while getting auth token from session", http.StatusInternalServerError)
					return
				}
				logger.Debug().Msgf("got auth token from session")
				r.Header.Set("Authorization", authToken)
				logger.Debug().Msgf("got auth token set to request header")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AuthFromSession(session map[interface{}]interface{}) (string, error) {
	var token interface{}
	var ok bool
	if token, ok = session[api.SessionTokenKey]; !ok {
		return "", errors.New("no token set on session")
	}
	if authToken, ok := token.(string); ok {
		return authToken, nil
	}
	return "", errors.New("value on token value is not string")
}
