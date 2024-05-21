package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const deckSessionCollection = "deck_sessions"

type (
	SessionDataAccess interface {
		GetActiveSessionForUserDeck(ctx context.Context, username string, deckID string) (models.DeckSession, error)
		GetSessionByID(ctx context.Context, sessionID string) (models.DeckSession, error)
		CreateSessionForUserDeck(ctx context.Context, session models.DeckSession) error
		UpdateCurrentCard(ctx context.Context, sessionID, currentCardID string, isFront bool) error
		SetAnswerForCard(ctx context.Context, sessionID, cardID string, isAnsweredCorrectly bool) error
		UpdateCardOrientation(ctx context.Context, sessionID string, isFront bool) error
		EndSession(ctx context.Context, sessionID string) error
	}
	SessionDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

// NewSessionDataAccess ...
func NewSessionDataAccess(db *mongo.Database, log zerolog.Logger) *SessionDAO {
	logger := log.With().Str("module", "SessionDAO").Logger()

	collection := db.Collection(deckSessionCollection)
	return &SessionDAO{
		collection: collection,
		log:        logger,
	}
}

func (s *SessionDAO) GetSessionByID(ctx context.Context, sessionID string) (models.DeckSession, error) {
	log := s.log.With().Str("method", "GetActiveSessionForUserDeck").Logger()

	filter := bson.D{
		{"_id", sessionID},
	}

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return models.DeckSession{}, ErrNoResults
		}
		log.Error().Err(result.Err()).Msgf("while looking up session for id %s", sessionID)
		return models.DeckSession{}, errors.Join(result.Err(), ErrFind)
	}

	var session = models.DeckSession{}
	err := result.Decode(&session)
	if err != nil {
		log.Error().Err(result.Err()).Msgf("while decoding session for id %s", sessionID)
		return models.DeckSession{}, result.Err()
	}
	return session, nil
}

func (s *SessionDAO) GetActiveSessionForUserDeck(ctx context.Context, username string, deckID string) (models.DeckSession, error) {
	log := s.log.With().Str("method", "GetActiveSessionForUserDeck").Logger()
	log.Info().Msgf("getting active session for user %s deck %s", username, deckID)

	filter := bson.D{
		{"deck_id", deckID},
		{"username", username},
		{"finished_at", nil},
	}

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return models.DeckSession{}, ErrNoResults
		}
		log.Error().Err(result.Err()).Msgf("while looking up session for username %s and deck id %s", username, deckID)
		return models.DeckSession{}, errors.Join(result.Err(), ErrFind)
	}

	var session = models.DeckSession{}
	err := result.Decode(&session)
	if err != nil {
		log.Error().Err(err).Msgf("while decoding session for username %s and deck id %s", username, deckID)
		return models.DeckSession{}, err
	}
	log.Debug().Msgf("results: %+v", session)
	return session, nil
}

func (s *SessionDAO) CreateSessionForUserDeck(ctx context.Context, session models.DeckSession) error {
	now := time.Now()
	log := s.log.With().Str("method", "CreateSessionForUserDeck").Logger()
	log.Info().Msgf("creating session for user %s deck %s", session.Username, session.DeckID)
	session.CreatedAt = now
	session.UpdatedAt = now
	_, err := s.collection.InsertOne(ctx, session)
	if err != nil {
		log.Error().Err(err).Msgf("while inserting session for %v", session)
		return errors.Join(err, ErrInsert)
	}
	log.Debug().Msgf("session created for user %s deck %s", session.Username, session.DeckID)
	return nil
}

func (s *SessionDAO) UpdateCurrentCard(ctx context.Context, sessionID, currentCardID string, isFront bool) error {
	log := s.log.With().Str("method", "UpdateCurrentCard").Logger()

	_, err := s.collection.UpdateOne(ctx, bson.D{
		{"_id", sessionID}},
		bson.D{
			{"$set", bson.D{
				{"current_card_id", currentCardID},
				{"is_front", isFront},
				{"updated_at", time.Now()},
			}},
		})
	if err != nil {
		log.Error().Err(err).Msgf("while updating session %s", sessionID)
		return errors.Join(err, ErrUpdate)
	}

	return nil
}

func (s *SessionDAO) SetAnswerForCard(ctx context.Context, sessionID, cardID string, isAnsweredCorrectly bool) error {
	log := s.log.With().Str("method", "SetAnswerForCard").Logger()

	filter := bson.D{
		{"_id", sessionID},
		{"card_answers.card_id", cardID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"updated_at", time.Now()},
			{"card_answers.$.card_id", cardID},
			{"card_answers.$.is_correct", isAnsweredCorrectly},
		}},
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msgf("while updating session %s", sessionID)
		return errors.Join(err, ErrUpdate)
	}

	if res.MatchedCount == 0 {
		filter = bson.D{
			{"_id", sessionID},
		}

		update = bson.D{
			{"$addToSet", bson.D{
				{"card_answers", bson.D{
					{"card_id", cardID},
					{"is_correct", isAnsweredCorrectly},
					{"created_at", time.Now()},
					{"updated_at", time.Now()},
				}},
			}},
			{"$set", bson.D{
				{"updated_at", time.Now()},
			}},
		}
		res, err = s.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Error().Err(err).Msgf("while updating session %s", sessionID)
			return errors.Join(err, ErrUpdate)
		}
		if res.MatchedCount == 0 {
			return errors.Join(ErrFind, fmt.Errorf("no match for session %s", sessionID))
		}
	}

	return nil
}

func (s *SessionDAO) UpdateCardOrientation(ctx context.Context, sessionID string, isFront bool) error {
	log := s.log.With().Str("method", "UpdateCardOrientation").Logger()

	_, err := s.collection.UpdateOne(ctx, bson.D{{"_id", sessionID}},
		bson.D{
			{"$set", bson.D{
				{"is_front", isFront},
				{"updated_at", time.Now()},
			}},
		})
	if err != nil {
		log.Error().Err(err).Msgf("while updating session %s", sessionID)
		return errors.Join(err, ErrUpdate)
	}

	return nil
}

func (s *SessionDAO) EndSession(ctx context.Context, sessionID string) error {
	log := s.log.With().Str("method", "EndSession").Logger()
	_, err := s.collection.UpdateOne(ctx, bson.D{{"_id", sessionID}},
		bson.D{
			{"$set", bson.D{
				{"finished_at", time.Now()},
			}},
		})
	log.Debug().Msg("called")
	if err != nil {
		log.Error().Err(err).Msgf("while ending session %s", sessionID)
		return errors.Join(err, ErrUpdate)
	}

	return nil
}
