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

var _ DeckDataAccess = &DeckDAO{}

type (
	DeckDataAccess interface {
		InsertDeck(ctx context.Context, deck models.Deck) (string, error)
		GetWithCards(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.DeckWithCards, error)
		GetDeckWithCardsByID(ctx context.Context, deckID string) (models.DeckWithCards, error)
		AddUserToUpvoteForDeck(ctx context.Context, primaryKey, userID string) error
		RemoveUserFromUpvoteForDeck(ctx context.Context, primaryKey, userID string) error
		AddUserToDownvoteForDeck(ctx context.Context, primaryKey, userID string) error
		RemoveUserFromDownvoteForDeck(ctx context.Context, primaryKey, userID string) error
		GetDecksForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.GetDeckResults, error)
	}

	DeckDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewDeckDataAccess(db *mongo.Database, log zerolog.Logger) *DeckDAO {
	logger := log.With().Str("module", "DataAccessObject").Logger()
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

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted deck")
		return "", errors.Join(fmt.Errorf("error inserting deck: %w", err), ErrInsert)
	}

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

func (d *DeckDAO) AddUserToUpvoteForDeck(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "AddUserToUpvoteForDeck").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"user_upvotes", userID},
		}},
		{"$pull", bson.D{
			{"user_downvotes", userID},
		}},
	}

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error adding user to upvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) RemoveUserFromUpvoteForDeck(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "RemoveUserFromUpvoteForDeck").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$pull", bson.D{
			{"user_upvotes", userID},
		}},
	}

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error removing user from upvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) AddUserToDownvoteForDeck(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "AddUserToDownvoteForDeck").Logger()
	logger.Info().Msgf("adding downvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"user_downvotes", userID},
		}},
		{"$pull", bson.D{
			{"user_upvotes", userID},
		}},
	}

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error adding user to downvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) RemoveUserFromDownvoteForDeck(ctx context.Context, deckID, userID string) error {
	logger := d.log.With().Str("method", "RemoveUserFromDownvoteForDeck").Logger()
	logger.Info().Msgf("adding upvote for user: %s", userID)

	filter := bson.D{{Key: "_id", Value: deckID}}
	update := bson.D{
		{"$pull", bson.D{
			{"user_upvotes", userID},
		}},
	}

	_, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("error removing user from downvote: %w", err), ErrUpdate)
	}

	return nil
}

func (d *DeckDAO) GetDeckWithCardsByID(ctx context.Context, deckID string) (models.DeckWithCards, error) {
	logger := d.log.With().Str("method", "updateCards").Logger()
	match := bson.D{{"$match", bson.D{{"_id", deckID}}}}
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
	sort := bson.D{{"$sort", bson.D{{"cards.created_at", 1}}}}

	p := bson.A{
		match,
		lookupCards,
		sort,
	}

	c, err := d.collection.Aggregate(ctx, p)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting cards on deckID of %s", deckID)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.DeckWithCards{}, ErrNoResults
		}
		return models.DeckWithCards{}, errors.Join(err, ErrAggregate)
	}
	defer c.Close(ctx)

	var decks []models.DeckWithCards
	err = c.All(ctx, &decks)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting deck from cursor for deckID: %s", deckID)
		return models.DeckWithCards{}, err
	}

	if len(decks) == 0 {
		return models.DeckWithCards{}, ErrNoResults
	}

	return decks[0], nil
}

func (d *DeckDAO) GetDecksForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.GetDeckResults, error) {
	logger := d.log.With().Str("method", "GetDecksForUser").Logger()
	logger.Info().Msgf("getting decks for user: %s", username)

	filter := []bson.D{
		bson.D{{"$match", bson.D{{"created_by", username}}}},
		bson.D{
			{"$addFields",
				bson.D{
					{"user_upvotes",
						bson.D{
							{"$ifNull",
								bson.A{
									"$user_upvotes",
									bson.A{},
								},
							},
						},
					},
					{"user_downvotes",
						bson.D{
							{"$ifNull",
								bson.A{
									"$user_downvotes",
									bson.A{},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"upvotes", bson.D{{"$size", "$user_upvotes"}}},
					{"downvotes", bson.D{{"$size", "$user_downvotes"}}},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cards"},
					{"localField", "_id"},
					{"foreignField", "deck_id"},
					{"as", "cards"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"name", "$name"},
					{"upvotes", "$upvotes"},
					{"downvotes", "$downvotes"},
					{"created_at", "$created_at"},
					{"created_updated", "$updated_at"},
					{"created_by", "$created_by"},
					{"num_cards", bson.D{
						{"$size",
							bson.D{
								{"$ifNull",
									bson.A{
										"$cards",
										bson.A{},
									},
								},
							},
						},
					}},
				},
			},
		},
	}

	filter = append(filter,
		pipeline.SortBy(pipeline.Asc))

	cur, err := d.collection.Aggregate(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Msg("error in calling aggregation")
		return nil, errors.Join(err, ErrAggregate)
	}
	defer cur.Close(ctx)

	var deckResults []models.GetDeckResults
	err = cur.All(ctx, &deckResults)
	if err != nil {
		logger.Error().Err(err).Msg("while unmarshalling into slice")
		return nil, errors.Join(err, ErrAggregate)
	}

	if len(deckResults) == 0 {
		logger.Info().Msgf("No decks belonging to %s", username)
		return nil, ErrNoResults
	}

	return deckResults, nil
}
