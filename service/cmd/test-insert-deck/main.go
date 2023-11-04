package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/internal/database/card"
	"github.com/rmarken/reptr/internal/database/deck"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
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

	deckDataAccess := deck.NewDataAccess(db, log)
	cardDataAccess := card.NewDataAccess(db, log)
	deckID, err := deckDataAccess.InsertDeck(ctx, deck.Deck{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Error().Err(err).Msg("Inserting deck")
	}

	err = cardDataAccess.InsertCards(ctx, []card.Card{
		{
			ID:        uuid.NewString(),
			Front:     "The host of Jeopardy",
			Back:      "Who is Alex Trebek",
			Kind:      card.BasicCard,
			DeckID:    deckID,
			CreatedAt: time.Now(),
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Inserting deck")
	}

}
