package cmd

import (
	"context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/decks/session"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	PORT              string `yaml:"PORT"`
	MongoUri          string `yaml:"MONGO_URI"`
	DbName            string `yaml:"DB_NAME"`
	Auth0Audience     string `yaml:"AUTH0_AUDIENCE"`
	Auth0ClientId     string `yaml:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `yaml:"AUTH0_CLIENT_SECRET"`
	Auth0GrantType    string `yaml:"AUTH0_GRANT_TYPE"`
	Auth0Endpoint     string `yaml:"AUTH0_ENDPOINT"`
	Auth0CallbackUrl  string `yaml:"AUTH0_CALLBACK_URL"`
	SessionKey        string `yaml:"SESSION_KEY"`
}

func LoadConfigFromFile(logger zerolog.Logger, path string) Config {
	file, err := os.ReadFile(path)
	if err != nil {
		logger.Panic().Err(err).Msg("while reading config file")
	}
	var conf Config
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		logger.Panic().Err(err).Msg("while marshalling config")
	}
	return conf
}

func LoadConfFromEnv(logger zerolog.Logger) Config {
	port := os.Getenv("PORT")
	if port == "" {
		logger.Panic().Msg("unable to get value for port")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if port == "" {
		logger.Panic().Msg("unable to get value for mongo")
	}
	dbName := os.Getenv("DB_NAME")
	if port == "" {
		logger.Panic().Msg("unable to get value for db name")
	}
	audience := os.Getenv("AUTH0_AUDIENCE")
	if audience == "" {
		logger.Panic().Msg("unable to get value for audience")
	}
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
	sessionKey := os.Getenv("SESSION_KEY")
	if callbackURL == "" {
		logger.Panic().Msg("unable to get value for session key")
	}
	config := Config{
		PORT:              port,
		MongoUri:          mongoURI,
		DbName:            dbName,
		Auth0Audience:     audience,
		Auth0ClientId:     clientID,
		Auth0ClientSecret: clientSecret,
		Auth0GrantType:    grantType,
		Auth0Endpoint:     authEndpoint,
		Auth0CallbackUrl:  callbackURL,
		SessionKey:        sessionKey,
	}
	return config
}

func MustLoadLogic(logger zerolog.Logger, repo database.Repository) *decks.Logic {
	return decks.New(logger, repo)
}

func MustLoadSessionLogic(logger zerolog.Logger, deckController decks.Controller, repo database.Repository) *session.Logic {
	return session.NewLogic(logger, deckController, repo)
}

func MustLoadProvider(logger zerolog.Logger, repo database.Repository) *provider.Logic {
	return provider.New(logger, repo)
}

func MustLoadRepo(logger zerolog.Logger, db *mongo.Database) *database.DataAccessObject {
	return database.NewRepository(logger, db)
}
func MustConnectMongo(ctx context.Context, logger zerolog.Logger, config Config) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoUri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Panic().Err(err).Msg("while connecting to mongo")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Panic().Err(err).Msg("while pinging mongo")
	}

	return client.Database(config.DbName)
}

func MustLoadAuth(ctx context.Context, logger zerolog.Logger, config Config, repo database.Repository) *auth.Authenticator {

	authenticator, err := auth.New(ctx, logger, repo, config.Auth0Audience, config.Auth0Endpoint, config.Auth0ClientId, config.Auth0ClientSecret, config.Auth0CallbackUrl)
	if err != nil {
		panic(err)
	}

	return authenticator
}
