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
	serverImpl := api.New(log, l, p, authenticator)

	router := mux.NewRouter()

	wrapper := exAPI.ServerInterfaceWrapper{
		Handler: serverImpl,
	}

	router.HandleFunc("/api/v1/login", wrapper.Login).Methods("GET")

	secureRoute := router.PathPrefix("/secure").Subrouter()
	secureRoute.HandleFunc("/api/v1/deck", wrapper.AddDeck).Methods("POST")
	secureRoute.HandleFunc("/api/v1/deck", wrapper.AddDeck).Methods("POST")
	secureRoute.HandleFunc("/api/v1/group", wrapper.AddGroup).Methods("POST")
	secureRoute.HandleFunc("/api/v1/group/{group_id}/deck/{deck_id}", wrapper.AddDeckToGroup).Methods("PUT")
	secureRoute.HandleFunc("/api/v1/groups", wrapper.GetGroups).Methods("GET")
	secureRoute.Use(middlewares.Authenticate(log, *authenticator))

	s := &http.Server{
		Handler: router,
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
