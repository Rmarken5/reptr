package main

import (
	"context"
	"github.com/rmarken/reptr/internal/database/deck"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const uri = "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Panic().Err(err).Msg("while connecting to mongo")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("deck")

	access := deck.NewDataAccess(db, log)
	err = access.InsertDeck(ctx, deck.Deck{
		Name: "my-test-deck",
	})

	if err != nil {
		log.Error().Err(err).Msg("Inserting deck")
	}
}
