package deck

import (
	"context"
	"errors"
	"github.com/rmarken/reptr/internal/database/pipeline"
	"github.com/rmarken/reptr/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	DataAccess interface {
		InsertDeck(ctx context.Context, deck models.Deck) (string, error)
		GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error)
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

func (d *DAO) InsertDeck(ctx context.Context, deck models.Deck) (string, error) {
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

func (d *DAO) GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error) {
	logger := d.log.With().Str("method", "GetWithCards").Logger()
	logger.Info().Msgf("Getting WithCards %v - %v, limit: %d offset %d", from, to, limit, offset)

	cur, err := d.collection.Find(ctx, pipeline.Paginate(from, to, limit, offset))
	if err != nil {
		return nil, err
	}
	withCards := make([]models.WithCards, 0)
	err = cur.All(ctx, &withCards)
	if err != nil {
		return nil, err
	}
	return withCards, nil
}
