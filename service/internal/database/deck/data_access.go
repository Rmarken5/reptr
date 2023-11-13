package deck

import (
	"context"
	"errors"
	"github.com/rmarken/reptr/internal/database/pipeline"
	"github.com/rmarken/reptr/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ DeckDataAccess = &DAO{}

type (
	DeckDataAccess interface {
		InsertDeck(ctx context.Context, deck models.Deck) (string, error)
		GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error)
		AddUserToUpvote(ctx context.Context, deckID, userID string) error
		RemoveUserFromUpvote(ctx context.Context, deckID, userID string) error
		AddUserToDownvote(ctx context.Context, deckID, userID string) error
		RemoveUserFromDownvote(ctx context.Context, deckID, userID string) error
	}

	DAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewDataAccess(db *mongo.Database, log zerolog.Logger) *DAO {
	logger := log.With().Str("module", "DAO").Logger()
	collection := db.Collection("decks")
	return &DAO{
		collection: collection,
		log:        logger,
	}
}

func (d *DAO) InsertDeck(ctx context.Context, deck models.Deck) (string, error) {
	logger := d.log.With().Str("method", "insertDeck").Logger()
	logger.Info().Msgf("inserting deck %v", deck)

	res, err := d.collection.InsertOne(ctx, deck)
	if err != nil {
		logger.Error().Err(err).Msgf("inserting deck %v", deck)
		return "", err
	}

	logger.Debug().Msgf("response: %+v", res.InsertedID)

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted deck")
		return "", errors.New("cannot return object id from inserted deck")
	}

	logger.Debug().Msgf("%s", prim)
	return prim, nil
}

func (d *DAO) GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error) {
	logger := d.log.With().Str("method", "GetWithCards").Logger()
	logger.Info().Msgf("Getting WithCards %v - %v, limit: %d offset %d", from, to, limit, offset)

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

	filter := append(
		pipeline.Paginate(from, to, limit, offset),
		lookupCards,
	)
	logger.Debug().Msgf("%+v", filter)

	cur, err := d.collection.Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}

	withCards := make([]models.WithCards, 0)
	err = cur.All(ctx, &withCards)
	if err != nil {
		return nil, err
	}
	return withCards, nil
}

func (d *DAO) AddUserToUpvote(ctx context.Context, deckID, userID string) error {
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
		return err
	}

	return nil
}

func (d *DAO) RemoveUserFromUpvote(ctx context.Context, deckID, userID string) error {
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
		return err
	}

	return nil
}

func (d *DAO) AddUserToDownvote(ctx context.Context, deckID, userID string) error {
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
		return err
	}

	return nil
}

func (d *DAO) RemoveUserFromDownvote(ctx context.Context, deckID, userID string) error {
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
		return err
	}

	return nil
}
