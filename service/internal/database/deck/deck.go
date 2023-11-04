package deck

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DataAccess interface {
		InsertDeck(ctx context.Context, deck Deck) (string, error)
	}

	DAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewDataAccess(db *mongo.Database, log zerolog.Logger) *DAO {
	logger := log.With().Str("module", "DAO").Logger()
	collection := db.Collection("decks")
	return &DAO{
		collection: collection,
		log:        logger,
	}
}

func (d *DAO) InsertDeck(ctx context.Context, deck Deck) (string, error) {
	logger := d.log.With().Str("method", "insertDeck").Logger()
	logger.Info().Msgf("inserting deck %v", deck)

	res, err := d.collection.InsertOne(ctx, deck)
	if err != nil {
		logger.Error().Err(err).Msgf("inserting deck %v", deck)
		return "", err
	}

	logger.Debug().Msgf("response: %+v", res.InsertedID)

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted deck")
		return "", errors.New("cannot return object id from inserted deck")
	}

	logger.Debug().Msgf("%s", prim)
	return prim, nil
}
