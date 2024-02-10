package session

import (
	"context"
	"errors"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
)

type (
	Controller interface {
		GetSessionForUserAndDeckID(ctx context.Context, username, deckID string) (models.DeckSession, error)
	}
	Logic struct {
		logger         zerolog.Logger
		repo           database.Repository
		deckController decks.Controller
	}
)

// NewLogic returns a pointer to a new Logic instance
func NewLogic(logger zerolog.Logger, deckController decks.Controller, repo database.Repository) *Logic {
	logger = logger.With().Str("module", "Session Logic").Logger()
	return &Logic{
		logger:         logger,
		repo:           repo,
		deckController: deckController,
	}
}

func (l Logic) GetSessionForUserAndDeckID(ctx context.Context, username, deckID string) (models.DeckSession, error) {
	log := l.logger.With().Str("method", "GetSessionForUserAndDeckID").Logger()
	log.Info().Msgf("getting deck session for username %s and deckID %s", username, deckID)
	var session models.DeckSession
	session, err := l.repo.GetSessionForUserDeck(ctx, deckID, username)
	if err != nil {
		if errors.Is(err, database.ErrNoResults) {
			deck, err := l.repo.GetDeckWithCardsByID(ctx, deckID)
			if err != nil {
				log.Error().Err(err).Msgf("while getting deckWithCards with ID %s", deckID)
				return models.DeckSession{}, err
			}
			currentCardID := ""
			if len(deck.Cards) > 0 {
				currentCardID = deck.Cards[0].ID
			}
			session = models.DeckSession{
				ID:            deckID + username,
				CurrentCardID: currentCardID,
				DeckName:      deck.Name,
				IsFront:       true,
			}
		}
	}
	return session, nil

}
