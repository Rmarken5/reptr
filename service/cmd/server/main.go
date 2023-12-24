package main

import (
	"context"
	"github.com/gorilla/mux"
	exAPI "github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/cmd"
	"github.com/rmarken/reptr/service/internal/api"
	"github.com/rmarken/reptr/service/internal/api/middlewares"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
)

const (
	osPort = "PORT"
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().Str("program", "reptr server").Logger()

	db := cmd.MustConnectMongo(ctx, log)
	defer db.Client().Disconnect(ctx)
	repo := cmd.MustLoadRepo(log, db)
	l := cmd.MustLoadLogic(log, repo)

	httpClient := http.Client{}
	authenticator, err := auth.New()
	if err != nil {
		log.Panic().Err(err).Msg("While creating authenticator")
	}
	p := cmd.MustLoadProvider(log, httpClient, repo)
	reptrClient := api.New(log, l, p, authenticator)

	r := mux.NewRouter()

	fromMux := exAPI.HandlerFromMux(reptrClient, r)
	secureRouter := r.PathPrefix("/secure/api/v1/groups").Subrouter()
	secureRouter.Use(middlewares.Authenticate(log, *authenticator))

	s := &http.Server{
		Handler: fromMux,
		Addr:    net.JoinHostPort("0.0.0.0", mustGetPort(log)),
	}

	log.Fatal().Err(s.ListenAndServe())

}

func mustGetPort(logger zerolog.Logger) string {
	port := os.Getenv(osPort)
	if port == "" {
		logger.Panic().Msgf("unable to get port")
	}
	return port
}
