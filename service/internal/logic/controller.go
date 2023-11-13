package logic

import (
	"context"
	"errors"
	"github.com/rmarken/reptr/internal/database"
	"github.com/rmarken/reptr/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	Controller interface {
		InsertDeck(ctx context.Context, deck models.Deck) (string, error)
		GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error)
		AddCardToDeck(ctx context.Context, deckID string, card models.Card) error
		UpdateCard(ctx context.Context, card models.Card) error
	}

	Logic struct {
		logger zerolog.Logger
		repo   database.Repository
	}
)

func New(logger zerolog.Logger, db *mongo.Database) *Logic {
	l := logger.With().Str("module", "logic").Logger()
	return &Logic{
		logger: l,
		repo:   database.NewRepository(l, db),
	}
}

// InsertDeck attempts to insert [models.Deck] into mongo. If repo returns an error, the error is logged and returned.
func (l *Logic) InsertDeck(ctx context.Context, deck models.Deck) (string, error) {
	logger := l.logger.With().Str("module", "InsertDeck").Logger()
	logger.Info().Msgf("insertDeck: %+v", deck)

	id, err := l.repo.InsertDeck(ctx, deck)
	if err != nil {
		l.logger.Error().Err(err).Msg("while inserting deck")
		return "", err
	}
	return id, nil
}

// GetDecks will attempt to get [[]models.WithCards] given a time period.
// From time is required. If to is not provided, it defaults to the EOD of from.
func (l *Logic) GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error) {
	logger := l.logger.With().Str("module", "getDecks").Logger()
	if to == nil {
		endOfDayFrom := time.Date(from.Year(), from.Month(), from.Day(), 23, 59, 59, 0, from.Location())
		to = &endOfDayFrom
	}
	logger.Info().Msgf("GetDecks between %s - %s with limit %d: starting at: %d ", from.Format(time.RFC3339), to.Format(time.RFC3339), limit, offset)

	if to.Before(from) {
		return nil, errors.New("to cannot be before from")
	}

	cards, err := l.repo.GetWithCards(ctx, from, to, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("while getting cards")
		return nil, err
	}

	return cards, nil
}

// AddCardToDeck adds a single [models.Card] to a given deck by ID.
func (l *Logic) AddCardToDeck(ctx context.Context, deckID string, card models.Card) error {
	logger := l.logger.With().Str("module", "addCardToDeck").Logger()
	logger.Info().Msgf("Adding card: %v to deck: %s", card, deckID)

	card.DeckID = deckID
	err := l.repo.InsertCards(ctx, []models.Card{card})
	if err != nil {
		logger.Error().Err(err).Msg("while inserting card")
		return err
	}

	return nil
}

// UpdateCard will update a card
func (l *Logic) UpdateCard(ctx context.Context, card models.Card) error {
	logger := l.logger.With().Str("module", "UpdateCard").Logger()
	logger.Info().Msgf("updating with card: %v", card)

	err := l.repo.UpdateCard(ctx, card)
	if err != nil {
		logger.Error().Err(err).Msgf("while updating card")
		return err
	}
	return nil
}
