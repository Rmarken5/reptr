package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

type (
	UserDataAccess interface {
		InsertUser(ctx context.Context, user models.User) (string, error)
		GetUserByID(ctx context.Context, id string) (models.User, error)
	}

	UserDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewUserDataAccess(db *mongo.Database, log zerolog.Logger) *UserDAO {
	logger := log.With().Str("module", "UserDAO").Logger()
	collection := db.Collection(usersCollection)
	return &UserDAO{
		collection: collection,
		log:        logger,
	}
}

func (u *UserDAO) InsertUser(ctx context.Context, user models.User) (string, error) {
	logger := u.log.With().Str("method", "insertUser").Logger()
	logger.Info().Msgf("Inserting user: %v", user)

	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting user %v", user)
		return "", errors.Join(err, ErrInsert)
	}

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted user")
		return "", errors.Join(fmt.Errorf("error inserting user: %w", err), ErrInsert)
	}

	logger.Debug().Msgf("%s", prim)

	return prim, nil
}

func (u *UserDAO) GetUserByID(ctx context.Context, id string) (models.User, error) {
	logger := u.log.With().Str("method", "getUserByID").Logger()

	one := u.collection.FindOne(ctx, bson.D{{"_id", id}})
	if one.Err() != nil {
		logger.Error().Err(one.Err()).Msgf("getting user by id %s", id)
		if errors.Is(one.Err(), mongo.ErrNoDocuments) {
			return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", one.Err()), ErrNoResults)
		}
		return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", one.Err()), ErrFind)
	}

	var usr models.User
	err := one.Decode(&usr)
	if err != nil {
		logger.Error().Err(err).Msgf("decoding result to user %s", id)
		return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", err))
	}

	return usr, nil
}
