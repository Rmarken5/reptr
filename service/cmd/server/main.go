package main

import (
	"context"
	"github.com/gorilla/mux"
	exAPI "github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/cmd"
	"github.com/rmarken/reptr/service/internal/api"
	"github.com/rmarken/reptr/service/internal/api/middlewares"
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

	authenticator := cmd.MustLoadAuth(ctx, log, repo)
	p := cmd.MustLoadProvider(log, repo)
	serverImpl := api.New(log, l, p, authenticator)

	router := mux.NewRouter()

	wrapper := exAPI.ServerInterfaceWrapper{
		Handler: serverImpl,
	}

	router.HandleFunc("/login", wrapper.LoginPage).Methods(http.MethodGet)
	router.HandleFunc("/login", wrapper.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", wrapper.Register).Methods(http.MethodPost)
	router.HandleFunc("/register", wrapper.RegistrationPage).Methods(http.MethodGet)

	secureRoute := router.PathPrefix("/secure").Subrouter()
	secureRoute.HandleFunc("/api/v1/deck", wrapper.AddDeck).Methods(http.MethodPost)
	secureRoute.HandleFunc("/api/v1/deck", wrapper.AddDeck).Methods(http.MethodPost)
	secureRoute.HandleFunc("/api/v1/group", wrapper.AddGroup).Methods(http.MethodPost)
	secureRoute.HandleFunc("/api/v1/group/{group_id}/deck/{deck_id}", wrapper.AddDeckToGroup).Methods("PUT")
	secureRoute.HandleFunc("/api/v1/groups", wrapper.GetGroups).Methods(http.MethodGet)
	secureRoute.Use(middlewares.Authenticate(log, authenticator))

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
