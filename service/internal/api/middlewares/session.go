package middlewares

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/service/internal/api"
	"github.com/rs/zerolog"
	"net/http"
)

//go:generate mockgen -destination ./mocks/sessions_mock.go -package middlewares  github.com/gorilla/sessions Store

func Session(logger zerolog.Logger, store sessions.Store) func(next http.Handler) http.Handler {
	logger = logger.With().Str("middleware", "Session").Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, api.CookieSessionID)
			if !session.IsNew {
				authToken, err := AuthFromSession(session.Values)
				if err != nil {
					logger.Error().Err(err).Msg("while getting auth token from session")
					http.Error(w, "while getting auth token from session", http.StatusBadRequest)
					return
				}
				r.Header.Set("Authorization", authToken)
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
