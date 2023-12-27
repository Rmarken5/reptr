package cmd

import (
	"context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
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

func MustLoadProvider(logger zerolog.Logger, repo database.Repository) *provider.Logic {
	return provider.New(logger, repo)
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

func MustLoadAuth(ctx context.Context, logger zerolog.Logger) *auth.Authenticator {
	//audience := os.Getenv("AUTH0_AUDIENCE")
	//if audience == "" {
	//	logger.Panic().Msg("unable to get value for audience")
	//}
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	if clientID == "" {
		logger.Panic().Msg("unable to get value for client id")
	}
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	if clientSecret == "" {
		logger.Panic().Msg("unable to get value for client secret")
	}
	grantType := os.Getenv("AUTH0_GRANT_TYPE")
	if grantType == "" {
		logger.Panic().Msg("unable to get value for grant type")
	}
	authEndpoint := os.Getenv("AUTH0_ENDPOINT")
	if authEndpoint == "" {
		logger.Panic().Msg("unable to get value for endpoint")
	}
	callbackURL := os.Getenv("AUTH0_CALLBACK_URL")
	if callbackURL == "" {
		logger.Panic().Msg("unable to get value for callbackURL")
	}
	authenticator, err := auth.New(ctx, authEndpoint, clientID, clientSecret, callbackURL)
	if err != nil {
		panic(err)
	}

	return authenticator
}
