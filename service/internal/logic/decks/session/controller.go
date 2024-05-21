package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type (
	Controller interface {
		GetActiveSessionForUserAndDeckID(ctx context.Context, username, deckID string) (models.DeckSession, error)
		UpdateSessionState(ctx context.Context, update models.SessionUpdate) error
		UpdateCardOrientation(ctx context.Context, sessionID string, isFront bool) error
		GetSessionByID(ctx context.Context, sessionID string) (models.DeckSession, error)
		SetCurrentCard(ctx context.Context, sessionID, cardID string, isFront bool) error
	}
	Logic struct {
		logger         zerolog.Logger
		repo           database.Repository
		deckController decks.Controller
	}
)

var _ Controller = &Logic{}

// NewLogic returns a pointer to a new Logic instance
func NewLogic(logger zerolog.Logger, deckController decks.Controller, repo database.Repository) *Logic {
	logger = logger.With().Str("module", "Session Logic").Logger()
	return &Logic{
		logger:         logger,
		repo:           repo,
		deckController: deckController,
	}
}

func (l *Logic) GetActiveSessionForUserAndDeckID(ctx context.Context, username, deckID string) (models.DeckSession, error) {
	log := l.logger.With().Str("method", "GetSessionForUserAndDeckID").Logger()
	log.Info().Msgf("getting deck session for username %s and deckID %s", username, deckID)
	var session models.DeckSession
	session, err := l.repo.GetActiveSessionForUserDeck(ctx, username, deckID)
	if err != nil {
		if errors.Is(err, database.ErrNoResults) {
			log.Debug().Msgf("no session for username %s and deckID %s", username, deckID)
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
				ID:            uuid.NewString(),
				Username:      username,
				DeckID:        deckID,
				DeckName:      deck.Name,
				CurrentCardID: currentCardID,
				IsFront:       true,
				CardAnswers:   make([]models.CardAnswer, 0),
			}

			sessionErr := l.repo.CreateSessionForUserDeck(ctx, session)
			if sessionErr != nil {
				log.Error().Err(sessionErr).Msgf("while creating session for user %s deck %s", username, deckID)
			}
		}
	}
	return session, nil
}

func (l *Logic) UpdateSessionState(ctx context.Context, update models.SessionUpdate) error {
	log := l.logger.With().Str("method", "UpdateSessionState").Logger()
	log.Info().Msgf("updating session state for session %s ", update.ID)

	err := l.repo.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var (
			g, errCtx = errgroup.WithContext(sessionContext)
		)
		g.Go(func() error {
			return l.repo.SetAnswerForCard(errCtx, update.ID, update.CurrentCardID, update.IsAnsweredCorrect)
		})
		g.Go(func() error {
			return l.repo.UpdateCurrentCard(errCtx, update.ID, update.NewCardID, update.IsFront)
		})
		if update.IsLastCard {
			g.Go(func() error {
				return l.repo.EndSession(errCtx, update.ID)
			})
		}
		if err := g.Wait(); err != nil {
			log.Error().Err(err).Msgf("while updating session state: %s", update.ID)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *Logic) UpdateCardOrientation(ctx context.Context, sessionID string, isFront bool) error {
	log := l.logger.With().Str("method", "UpdateCardOrientation").Logger()
	log.Info().Msgf("updating card orientation for session %s", sessionID)

	return l.repo.UpdateCardOrientation(ctx, sessionID, isFront)
}

func (l *Logic) GetSessionByID(ctx context.Context, sessionID string) (models.DeckSession, error) {
	log := l.logger.With().Str("method", "GetSessionByID").Logger()
	log.Info().Msgf("getting session by id %s", sessionID)

	return l.repo.GetSessionByID(ctx, sessionID)
}

func (l *Logic) SetCurrentCard(ctx context.Context, sessionID string, cardID string, isFront bool) error {
	log := l.logger.With().Str("method", "SetCurrentCard").Logger()
	log.Info().Msgf("setting current card for session %s", sessionID)
	return l.repo.UpdateCurrentCard(ctx, sessionID, cardID, isFront)
}
