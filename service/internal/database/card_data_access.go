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

var _ CardDataAccess = &CardDAO{}

type (
	CardDataAccess interface {
		InsertCards(ctx context.Context, card []models.Card) error
		UpdateCard(ctx context.Context, card models.Card) error
	}
	CardDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewCardDataAccess(db *mongo.Database, log zerolog.Logger) *CardDAO {
	logger := log.With().Str("module", "CardDAO").Logger()
	collection := db.Collection("cards")
	return &CardDAO{
		collection: collection,
		log:        logger,
	}
}

func (d *CardDAO) InsertCards(ctx context.Context, cards []models.Card) error {
	logger := d.log.With().Str("method", "insertCard").Logger()
	logger.Info().Msgf("Inserting cards %v", cards)

	c := make([]interface{}, 0, len(cards))
	for _, card := range cards {
		c = append(c, card)
	}

	_, err := d.collection.InsertMany(ctx, c)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting cards %v", cards)
		return errors.Join(fmt.Errorf("error inserting cards: %w", err), ErrInsert)
	}

	return nil
}

func (d *CardDAO) UpdateCard(ctx context.Context, card models.Card) error {
	logger := d.log.With().Str("method", "updateCards").Logger()

	filter := bson.D{{"_id", card.ID}}

	u, err := d.collection.UpdateOne(ctx, filter, bson.M{
		"$set": card,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("Updating card %v", card)
		return errors.Join(fmt.Errorf("error updating card: %w", err), ErrUpdate)
	}

	logger.Info().Msgf("Updated: %+v", u)

	return nil
}
