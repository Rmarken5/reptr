package logic

import (
	"context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ Controller = &Logic{}

type (
	Controller interface {
		CreateGroup(ctx context.Context, group models.Group) (string, error)
		AddDeckToGroup(ctx context.Context, groupID, deckID string) error
		GetGroups(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error)
		CreateDeck(ctx context.Context, deck models.Deck) (string, error)
		GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.DeckWithCards, error)
		AddCardToDeck(ctx context.Context, deckID string, card models.Card) error
		UpdateCard(ctx context.Context, card models.Card) error
		UpvoteDeck(ctx context.Context, deckID, userID string) error
		RemoveUpvoteDeck(ctx context.Context, deckID, userID string) error
		DownvoteDeck(ctx context.Context, deckID, userID string) error
		RemoveDownvoteDeck(ctx context.Context, deckID, userID string) error
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

// CreateDeck attempts to insert [models.Deck] into mongo. If repo returns an error, the error is logged and returned.
func (l *Logic) CreateDeck(ctx context.Context, deck models.Deck) (string, error) {
	logger := l.logger.With().Str("module", "CreateDeck").Logger()
	logger.Info().Msgf("insertDeck: %+v", deck)

	id, err := l.repo.InsertDeck(ctx, deck)
	if err != nil {
		l.logger.Error().Err(err).Msg("while inserting deck")
		return "", err
	}
	return id, nil
}

// GetDecks will attempt to get [[]models.DeckWithCards] given a time period.
// From time is required. If to is not provided, it defaults to the EOD of from.
func (l *Logic) GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.DeckWithCards, error) {
	logger := l.logger.With().Str("module", "getDecks").Logger()
	if to == nil {
		endOfDayFrom := time.Date(from.Year(), from.Month(), from.Day(), 23, 59, 59, 0, from.Location())
		to = &endOfDayFrom
	}
	logger.Info().Msgf("GetDecks between %s - %s with limit %d: starting at: %d ", from.Format(time.RFC3339), to.Format(time.RFC3339), limit, offset)

	if to.Before(from) {
		return nil, ErrInvalidToBeforeFrom
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

func (l *Logic) UpvoteDeck(ctx context.Context, deckID, userID string) error {
	logger := l.logger.With().Str("module", "UpvoteDeck").Logger()
	logger.Info().Msgf("Upvote deck %s for %s", deckID, userID)

	err := l.repo.AddUserToUpvote(ctx, deckID, userID)
	if err != nil {
		logger.Error().Err(err).Msgf("while upvoting deck")
		return err
	}
	return nil
}

func (l *Logic) RemoveUpvoteDeck(ctx context.Context, deckID, userID string) error {
	logger := l.logger.With().Str("module", "RemoveUpvoteDeck").Logger()
	logger.Info().Msgf("remove upvote: deck %s for %s", deckID, userID)

	err := l.repo.RemoveUserFromUpvote(ctx, deckID, userID)
	if err != nil {
		logger.Error().Err(err).Msgf("while removing upvote")
		return err
	}
	return nil
}

func (l *Logic) DownvoteDeck(ctx context.Context, deckID, userID string) error {
	logger := l.logger.With().Str("module", "DownvoteDeck").Logger()
	logger.Info().Msgf("downvote deck: %s for %s", deckID, userID)

	err := l.repo.AddUserToDownvote(ctx, deckID, userID)
	if err != nil {
		logger.Error().Err(err).Msgf("while adding downvote")
		return err
	}
	return nil
}

func (l *Logic) RemoveDownvoteDeck(ctx context.Context, deckID, userID string) error {
	logger := l.logger.With().Str("module", "RemoveDownvoteDeck").Logger()
	logger.Info().Msgf("remove downvote deck: %s for %s", deckID, userID)

	err := l.repo.RemoveUserFromDownvote(ctx, deckID, userID)
	if err != nil {
		logger.Error().Err(err).Msgf("while removing downvote")
		return err
	}
	return nil
}

func (l *Logic) CreateGroup(ctx context.Context, group models.Group) (string, error) {
	logger := l.logger.With().Str("module", "CreateGroup").Logger()
	logger.Info().Msgf("CreateGroup: %+v", group)
	gpID, err := l.repo.InsertGroup(ctx, group)
	if err != nil {
		l.logger.Error().Err(err).Msg("while inserting group")
		return "", err
	}
	return gpID, nil
}

func (l *Logic) AddDeckToGroup(ctx context.Context, groupID, deckID string) error {
	logger := l.logger.With().Str("module", "AddDeckToGroup").Logger()
	logger.Info().Msgf("Adding deck: %s to group: %s", deckID, groupID)

	err := l.repo.AddDeckToGroup(ctx, groupID, deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while adding deck: %s to group: %s", deckID, groupID)
		return err
	}

	return nil
}

func (l *Logic) GetGroups(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error) {
	logger := l.logger.With().Str("module", "GetGroups").Logger()
	if to == nil {
		endOfDayFrom := time.Date(from.Year(), from.Month(), from.Day(), 23, 59, 59, 0, from.Location())
		to = &endOfDayFrom
	}
	logger.Info().Msgf("GetGroups between %s - %s with limit %d: starting at: %d ", from.Format(time.RFC3339), to.Format(time.RFC3339), limit, offset)

	if to.Before(from) {
		return nil, ErrInvalidToBeforeFrom
	}

	groupsWithDecks, err := l.repo.GetGroupsWithDecks(ctx, from, to, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("while getting groupsWithDecks")
		return nil, err
	}

	return groupsWithDecks, nil
}
