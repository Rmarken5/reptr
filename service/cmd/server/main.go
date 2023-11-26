package main

import (
	"context"
	"github.com/gorilla/mux"
	exAPI "github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/cmd"
	"github.com/rmarken/reptr/service/internal/api"
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

	l := cmd.MustLoadLogic(log, db)
	reptrClient := api.New(log, l)

	// This is how you set up a basic Gorilla router
	r := mux.NewRouter()

	// We now register our petStore above as the handler for the interface
	exAPI.HandlerFromMux(reptrClient, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", mustGetPort(log)),
	}

	// And we serve HTTP until the world ends.
	log.Fatal().Err(s.ListenAndServe())

}

func mustGetPort(logger zerolog.Logger) string {
	port := os.Getenv(osPort)
	if port == "" {
		logger.Panic().Msgf("unable to get port")
	}
	return port
}
