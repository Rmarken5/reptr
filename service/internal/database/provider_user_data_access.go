package database

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const providerUsersCollection = "provider_users"

type (
	ProviderUsersDataAccess interface {
		InsertUserSubjectPair(ctx context.Context, userID, subject string) error
		GetUserIDFor(ctx context.Context, subject string) (string, error)
	}

	ProviderUsersDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewProviderUsersDataAccess(db *mongo.Database, log zerolog.Logger) *ProviderUsersDAO {
	logger := log.With().Str("module", "ProviderUserDAO").Logger()
	collection := db.Collection(providerUsersCollection)
	return &ProviderUsersDAO{
		collection: collection,
		log:        logger,
	}
}

func (dao *ProviderUsersDAO) InsertUserSubjectPair(ctx context.Context, userID, subject string) error {
	logger := dao.log.With().Str("method", "InsertUserSubjectPair").Logger()
	logger.Info().Msgf("inserting user: %s subject: %s pair", userID, subject)

	_, err := dao.collection.InsertOne(ctx, bson.D{
		{"_id", subject},
		{"user_id", userID},
	})
	if err != nil {
		logger.Error().Err(err).Msgf("while inserting provider_id/user_id pair: %s/%s", subject, userID)
		return errors.Join(err, ErrInsert)
	}

	return nil
}
func (dao *ProviderUsersDAO) GetUserIDFor(ctx context.Context, subject string) (string, error) {
	var result struct {
		UserID string `bson:"user_id"`
	}
	err := dao.collection.FindOne(ctx, bson.M{"_id": subject}, options.FindOne().SetProjection(bson.M{"user_id": 1})).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrNoResults
		}
		return "", errors.Join(err, ErrFind)
	}

	return result.UserID, nil
}
