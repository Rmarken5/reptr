package database

import (
	"context"
	"errors"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SessionDataAccess interface {
		GetSessionForUserDeck(ctx context.Context, username string, deckID string) (models.DeckSession, error)
		CreateSessionForUserDeck(ctx context.Context, session models.DeckSession) error
		UpdateSession(ctx context.Context, username, deckID, currentCardID string, isFront bool) error
	}
	SessionDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

// NewSessionDataAccess ...
func NewSessionDataAccess(db *mongo.Database, log zerolog.Logger) *SessionDAO {
	logger := log.With().Str("module", "SessionDAO").Logger()
	collection := db.Collection("deck_sessions")
	return &SessionDAO{
		collection: collection,
		log:        logger,
	}
}

func (s SessionDAO) GetSessionForUserDeck(ctx context.Context, username string, deckID string) (models.DeckSession, error) {
	log := s.log.With().Str("method", "GetSessionForUserDeck").Logger()

	filter := bson.D{
		{"_id", deckID + username},
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
		log.Error().Err(result.Err()).Msgf("while decoding session for username %s and deck id %s", username, deckID)
		return models.DeckSession{}, result.Err()
	}
	return session, nil
}

func (s SessionDAO) CreateSessionForUserDeck(ctx context.Context, session models.DeckSession) error {
	log := s.log.With().Str("method", "CreateSessionForUserDeck").Logger()

	_, err := s.collection.InsertOne(ctx, session)
	if err != nil {
		log.Error().Err(err).Msgf("while inserting session for %v", session)
		return errors.Join(err, ErrInsert)
	}

	return nil
}

func (s SessionDAO) UpdateSession(ctx context.Context, username, deckID, currentCardID string, isFront bool) error {
	log := s.log.With().Str("method", "UpdateSession").Logger()

	_, err := s.collection.UpdateOne(ctx, bson.D{
		{"_id", deckID + username}},
		bson.D{
			{"$set", bson.D{
				{"current_card_id", currentCardID},
				{"is_front", isFront},
			}},
		})
	if err != nil {
		log.Error().Err(err).Msgf("while updating username %s and deckID %s", username, deckID)
		return errors.Join(err, ErrUpdate)
	}

	return nil
}
