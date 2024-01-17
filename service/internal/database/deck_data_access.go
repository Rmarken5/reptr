package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/service/internal/database/pipeline"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ DeckDataAccess = &DAO{}

type (
	DeckDataAccess interface {
		InsertDeck(ctx context.Context, deck models.Deck) (string, error)
		GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.DeckWithCards, error)
		AddUserToUpvote(ctx context.Context, deckID, userID string) error
		RemoveUserFromUpvote(ctx context.Context, deckID, userID string) error
		AddUserToDownvote(ctx context.Context, deckID, userID string) error
		RemoveUserFromDownvote(ctx context.Context, deckID, userID string) error
	}

	DeckDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewDeckDataAccess(db *mongo.Database, log zerolog.Logger) *DeckDAO {
	logger := log.With().Str("module", "DAO").Logger()
	collection := db.Collection("decks")
	return &DeckDAO{
		collection: collection,
		log:        logger,
	}
}

func (d *DeckDAO) InsertDeck(ctx context.Context, deck models.Deck) (string, error) {
	logger := d.log.With().Str("method", "insertDeck").Logger()
	logger.Info().Msgf("inserting deck %v", deck)

	res, err := d.collection.InsertOne(ctx, deck)
	if err != nil {
		logger.Error().Err(err).Msgf("inserting deck %v", deck)
		return "", errors.Join(err, ErrInsert)
	}

	logger.Debug().Msgf("response: %+v", res.InsertedID)

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted deck")
		return "", errors.Join(fmt.Errorf("error inserting deck: %w", err), ErrInsert)
	}

	logger.Debug().Msgf("%s", prim)
	return prim, nil
}

func (d *DeckDAO) GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.DeckWithCards, error) {
	logger := d.log.With().Str("method", "GetWithCards").Logger()
	logger.Info().Msgf("Getting DeckWithCards %v - %v, limit: %d offset %d", from, to, limit, offset)

	lookupCards :=
		bson.D{{"$lookup",
			bson.D{
				{"from", "cards"},
				{"localField", "_id"},
				{"foreignField", "deck_id"},
				{"as", "cards"},
			},
		},
		}

	getVotes := bson.D{
		{"$addFields",
			bson.D{
				{"upvotes", bson.D{{"$size", "$user_upvotes"}}},
				{"downvotes", bson.D{{"$size", "$user_downvotes"}}},
			},
		},
	}
	removeUserVotes := bson.D{
		{"$project",
			bson.D{
				{"user_downvotes", 0},
				{"user_upvotes", 0},
			},
		},
	}

	filter := append(
		pipeline.Paginate(from, to, limit, offset),
		getVotes,
		removeUserVotes,
		lookupCards,
	)
	logger.Debug().Msgf("%+v", filter)

	cur, err := d.collection.Aggregate(ctx, filter)
	if err != nil {
		return nil, errors.Join(err, ErrAggregate)
	}

	withCards := make([]models.DeckWithCards, 0)
	err = cur.All(ctx, &withCards)
	if err != nil {
		return nil, err
	}
	if len(withCards) == 0 {
		return nil, ErrNoResults
	}
	return withCards, nil
}

func (d *DeckDAO) AddUserToUpvote(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "AddUserToUpvote").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"user_upvote", userID},
		}},
		{"$pull", bson.D{
			{"user_downvote", userID},
		}},
	}

	logger.Debug().Msgf("%+v", filter)

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error adding user to upvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) RemoveUserFromUpvote(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "AddUserToUpvote").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$pull", bson.D{
			{"user_upvote", userID},
		}},
	}

	logger.Debug().Msgf("%+v", filter)

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error removing user from upvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) AddUserToDownvote(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "AddUserToDownvote").Logger()
	logger.Info().Msgf("adding downvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"user_downvote", userID},
		}},
		{"$pull", bson.D{
			{"user_upvote", userID},
		}},
	}

	logger.Debug().Msgf("%+v", filter)

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error adding user to downvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) RemoveUserFromDownvote(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "RemoveUserFromDownvote").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$pull", bson.D{
			{"user_upvote", userID},
		}},
	}

	logger.Debug().Msgf("%+v", filter)

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error removing user from downvote: %w", err), ErrUpdate)
	}

	return nil
}
