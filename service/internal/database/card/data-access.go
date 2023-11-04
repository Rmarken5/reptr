package card

import (
	"context"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DataAccess interface {
		InsertCards(ctx context.Context, card []Card) error
	}
	DAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewDataAccess(db *mongo.Database, log zerolog.Logger) *DAO {
	logger := log.With().Str("module", "CardDAO").Logger()
	collection := db.Collection("cards")
	return &DAO{
		collection: collection,
		log:        logger,
	}
}

func (d *DAO) InsertCards(ctx context.Context, cards []Card) error {
	logger := d.log.With().Str("method", "insertCard").Logger()
	logger.Info().Msgf("Inserting cards %v", cards)

	c := make([]interface{}, 0, len(cards))
	for _, card := range cards {
		c = append(c, card)
	}

	_, err := d.collection.InsertMany(ctx, c)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting cards %v", cards)
		return err
	}

	return nil
}
