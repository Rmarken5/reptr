package decks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"time"
)

//go:generate mockgen -destination ./mocks/controller_mock.go -package logic . Controller
var _ Controller = &Logic{}

type (
	Controller interface {
		CreateGroup(ctx context.Context, username, groupName string) (string, error)
		AddDeckToGroup(ctx context.Context, groupID, deckID string) error
		GetGroups(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error)

		GetGroupByID(ctx context.Context, groupID string) (models.GroupWithDecks, error)
		GetGroupsForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.Group, error)
		CreateDeck(ctx context.Context, deckName string) (string, error)
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

func New(logger zerolog.Logger, repo database.Repository) *Logic {
	l := logger.With().Str("module", "deck logic").Logger()
	return &Logic{
		logger: l,
		repo:   repo,
	}
}

// CreateDeck attempts to insert [models.Deck] into mongo. If repo returns an error, the error is logged and returned.
func (l *Logic) CreateDeck(ctx context.Context, deckName string) (string, error) {
	logger := l.logger.With().Str("module", "CreateDeck").Logger()
	logger.Info().Msgf("insertDeck: %s", deckName)

	if deckName == "" {
		logger.Error().Err(ErrEmptyDeckName).Msgf("deckName: %s", deckName)
		return "", ErrEmptyDeckName
	}

	id, err := l.repo.InsertDeck(ctx, models.Deck{
		ID:        uuid.NewString(),
		Name:      deckName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
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

	if to != nil && to.Before(from) {
		return []models.DeckWithCards(nil), ErrInvalidToBeforeFrom
	}

	cards, err := l.repo.GetWithCards(ctx, from, to, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("while getting cards")
		return []models.DeckWithCards(nil), err
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

func (l *Logic) CreateGroup(ctx context.Context, username, groupName string) (string, error) {
	logger := l.logger.With().Str("module", "CreateGroup").Logger()
	logger.Info().Msgf("CreateGroup: %s", groupName)

	if username == "" {
		logger.Error().Err(ErrEmptyUsername)
		return "", ErrEmptyUsername
	}

	if groupName == "" {
		logger.Error().Err(ErrInvalidGroupName)
		return "", ErrInvalidGroupName
	}
	gpID, err := l.repo.InsertGroup(ctx, models.Group{
		ID:         uuid.NewString(),
		Name:       groupName,
		CreatedBy:  username,
		Moderators: []string{username},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	})
	if err != nil {
		l.logger.Error().Err(err).Msg("while inserting group")
		return "", err
	}

	err = l.repo.AddUserAsMemberOfGroup(ctx, username, gpID)
	if err != nil {
		l.logger.Error().Err(err).Msgf("while making user %s member of group %s", username, groupName)
		return "", err
	}
	return gpID, nil
}

func (l *Logic) AddDeckToGroup(ctx context.Context, groupID, deckID string) error {
	logger := l.logger.With().Str("module", "AddDeckToGroup").Logger()
	logger.Info().Msgf("Adding deck: %s to group: %s", deckID, groupID)

	if groupID == "" {
		logger.Error().Err(ErrEmptyGroupID).Msgf("group: %s", groupID)
		return ErrEmptyGroupID
	}

	if deckID == "" {
		logger.Error().Err(ErrEmptyDeckID).Msgf("deck: %s", deckID)
		return ErrEmptyDeckID
	}

	err := l.repo.AddDeckToGroup(ctx, groupID, deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while adding deck: %s to group: %s", deckID, groupID)
		return err
	}

	return nil
}

func (l *Logic) GetGroups(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error) {
	logger := l.logger.With().Str("module", "GetGroups").Logger()

	if to != nil && to.Before(from) {
		return []models.GroupWithDecks(nil), ErrInvalidToBeforeFrom
	}

	groupsWithDecks, err := l.repo.GetGroupsWithDecks(ctx, from, to, limit, offset)
	if err != nil && !errors.Is(err, database.ErrNoResults) {
		logger.Error().Err(err).Msg("while getting groupsWithDecks")
		return []models.GroupWithDecks(nil), err
	}

	return groupsWithDecks, nil
}

func (l *Logic) GetGroupsForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.Group, error) {
	logger := l.logger.With().Str("method", "GetGroupsForUser").Logger()
	logger.Info().Msgf("GetGroupsForUser called")

	if username == "" {
		logger.Error().Err(ErrEmptyUsername)
		return []models.Group(nil), ErrEmptyUsername
	}

	if to != nil && to.Before(from) {
		return []models.Group(nil), ErrInvalidToBeforeFrom
	}

	groups, err := l.repo.GetGroupsForUser(ctx, username, from, to, limit, offset)
	if err != nil && !errors.Is(err, database.ErrNoResults) {
		logger.Error().Err(err).Msgf("while getting groups for user: %s", username)
		return []models.Group(nil), err
	}

	return groups, nil
}

func (l *Logic) GetGroupByID(ctx context.Context, groupID string) (models.GroupWithDecks, error) {
	logger := l.logger.With().Str("module", "GetGroupByID").Logger()

	group, err := l.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		logger.Error().Err(err).Msg("while getting cards")
		return models.GroupWithDecks{}, err
	}

	return group, nil
}
