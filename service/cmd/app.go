package cmd

import (
	"context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	mongoURI = "MONGO_URI"
	dbName   = "DB_NAME"
)

func MustLoadLogic(logger zerolog.Logger, repo database.Repository) *decks.Logic {
	return decks.New(logger, repo)
}

func MustLoadRepo(logger zerolog.Logger, db *mongo.Database) *database.DAO {
	return database.NewRepository(logger, db)
}
func MustConnectMongo(ctx context.Context, logger zerolog.Logger) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mustLoadMongoURI(logger)).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Panic().Err(err).Msg("while connecting to mongo")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Panic().Err(err).Msg("while pinging mongo")
	}

	return client.Database(mustLoadMongoDBName(logger))
}
func mustLoadMongoURI(logger zerolog.Logger) string {
	uri := os.Getenv(mongoURI)
	if uri == "" {
		logger.Panic().Msg("unable to get value for mongo uri")
	}
	return uri
}

func mustLoadMongoDBName(logger zerolog.Logger) string {
	dbName := os.Getenv(dbName)
	if dbName == "" {
		logger.Panic().Msg("unable to get value for database name")
	}
	return dbName
}
