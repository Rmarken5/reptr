package deck_viewer

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type (
	Controller interface {
		AnswerCardCorrect(ctx context.Context, sessionID string) (templ.Component, error)
	}

	Logic struct {
		logger zerolog.Logger
		repo   database.Repository
	}
)

func New(log zerolog.Logger, repo database.Repository) *Logic {
	logger = log.With().Str("service", "deck-viewer").Logger()
	return &Logic{
		logger: log,
		repo:   repo,
	}
}

func (l *Logic) AnswerCardCorrect(ctx context.Context, sessionID string) (templ.Component, error) {
	log := l.logger.With().Str("component", "AnswerCardCorrect").Logger()
	log.Info().Msgf("updating card correct for session: %s", sessionID)

	session, err := l.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		log.Error().Err(err).Msg("while getting session")
		return nil, err
	}

	frontOfCard, err := l.repo.GetFrontOfNextCardByID(ctx, session.DeckID, session.CurrentCardID, session.Username)
	if err != nil {
		// End of session
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = l.repo.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
				var (
					g, errCtx = errgroup.WithContext(sessionContext)
				)
				g.Go(func() error {
					return l.repo.SetAnswerForCard(errCtx, sessionID, session.CurrentCardID, true)
				})
				g.Go(func() error {
					return l.repo.UpdateCurrentCard(errCtx, sessionID, session.CurrentCardID, false)
				})

				g.Go(func() error {
					return l.repo.EndSession(errCtx, sessionID)
				})

				if err := g.Wait(); err != nil {
					log.Error().Err(err).Msgf("while updating session state: %s", sessionID)
					return nil, err
				}
				return nil, nil
			})
			if err != nil {
				return nil, err
			}
			// End of session

		}
		log.Error().Err(err).Msg("while getting front of card")
		return nil, err
	}

	err = l.repo.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var (
			g, errCtx = errgroup.WithContext(sessionContext)
		)
		g.Go(func() error {
			return l.repo.SetAnswerForCard(errCtx, sessionID, session.CurrentCardID, true)
		})
		g.Go(func() error {
			return l.repo.UpdateCurrentCard(errCtx, sessionID, frontOfCard.CardID, true)
		})

		if err := g.Wait(); err != nil {
			log.Error().Err(err).Msgf("while updating session state: %s", sessionID)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	// Return next card
	return

}
