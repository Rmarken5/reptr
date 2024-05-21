package deck_viewer

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/web/components/dumb"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

type (
	Controller interface {
		AnswerCurrentCard(ctx context.Context, sessionID string, isAnsweredCorrect bool) (templ.Component, error)
	}

	Logic struct {
		logger zerolog.Logger
		repo   database.Repository
	}
)

func New(log zerolog.Logger, repo database.Repository) *Logic {
	log = log.With().Str("service", "deck-viewer").Logger()
	return &Logic{
		logger: log,
		repo:   repo,
	}
}

func (l *Logic) AnswerCurrentCard(ctx context.Context, sessionID string, isAnsweredCorrect bool) (templ.Component, error) {
	log := l.logger.With().Str("component", "AnswerCurrentCard").Logger()
	log.Info().Msgf("updating card correct for session: %s", sessionID)

	session, err := l.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		log.Error().Err(err).Msg("while getting session")
		return nil, err
	}

	frontOfCard, err := l.repo.GetFrontOfNextCardByID(ctx, session.DeckID, session.CurrentCardID, session.Username)
	if err != nil {
		// End of session
		if errors.Is(err, database.ErrNoResults) {
			err = l.repo.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {

				err2 := l.repo.SetAnswerForCard(sessionContext, sessionID, session.CurrentCardID, isAnsweredCorrect)
				if err2 != nil {
					log.Error().Err(err2).Msg("while updating current card")
					return nil, err2
				}

				err2 = l.repo.UpdateCurrentCard(sessionContext, sessionID, session.CurrentCardID, false)
				if err2 != nil {
					log.Error().Err(err2).Msg("while updating current card")
					return nil, err2
				}

				err2 = l.repo.EndSession(sessionContext, sessionID)
				if err2 != nil {
					log.Error().Err(err2).Msg("while ending session")
					return nil, err2
				}
				return nil, nil
			})

			if err != nil {
				return nil, err
			}
			// End of session
			backOfCard, err := l.repo.GetBackOfCardByID(ctx, session.DeckID, session.CurrentCardID, session.Username)
			if err != nil {
				return dumb.BackOfCardDisplay(dumb.CardBack{}), nil
			}
			return dumb.BackOfCardDisplay(dumb.CardBack{
				SessionID:      sessionID,
				DeckID:         session.DeckID,
				CardID:         backOfCard.CardID,
				BackContent:    backOfCard.Answer,
				NextCardID:     backOfCard.NextCard,
				PreviousCardID: session.CurrentCardID,
				IsUpvoted:      bool(backOfCard.IsUpvotedByUser),
				IsDownvoted:    bool(backOfCard.IsDownvotedByUser),
				VoteButtonData: dumb.VoteButtonsData{
					CardID:            backOfCard.CardID,
					UpvoteClass:       backOfCard.IsUpvotedByUser.UpvotedClass(),
					DownvoteClass:     backOfCard.IsDownvotedByUser.DownvotedClass(),
					UpvoteDirection:   backOfCard.IsUpvotedByUser.NextUpvoteDirection(),
					DownvoteDirection: backOfCard.IsDownvotedByUser.DownvotedClass()},
			}), nil
		}
		log.Error().Err(err).Msg("while getting front of card")
		return nil, err
	}

	err = l.repo.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {

		err2 := l.repo.SetAnswerForCard(sessionContext, sessionID, session.CurrentCardID, isAnsweredCorrect)
		if err2 != nil {
			log.Error().Err(err2).Msg("while updating current card")
			return nil, err2
		}

		err2 = l.repo.UpdateCurrentCard(sessionContext, sessionID, frontOfCard.CardID, true)
		if err2 != nil {
			log.Error().Err(err2).Msg("while updating current card")
			return nil, err2
		}

		return nil, nil
	})
	if err != nil {
		log.Error().Err(err).Msg("while answering current card")
		return nil, err
	}

	// Return next card
	return dumb.FrontCardDisplay(dumb.CardFront{
		SessionID:      sessionID,
		DeckID:         session.DeckID,
		CardID:         frontOfCard.CardID,
		Front:          frontOfCard.Content,
		NextCardID:     frontOfCard.NextCard,
		PreviousCardID: session.CurrentCardID,
		Downvotes:      strconv.Itoa(frontOfCard.Downvotes),
		Upvotes:        strconv.Itoa(frontOfCard.Upvotes),
		CardType:       "",
	}), nil
}
