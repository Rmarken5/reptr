package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ CardDataAccess = &CardDAO{}

type (
	CardDataAccess interface {
		InsertCards(ctx context.Context, card []models.Card) error
		UpdateCard(ctx context.Context, card models.Card) error
		GetFrontOfCardByID(ctx context.Context, deckID, cardID string) (models.FrontOfCard, error)
		GetBackOfCardByID(ctx context.Context, deckID, cardID string) (models.BackOfCard, error)
	}
	CardDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewCardDataAccess(db *mongo.Database, log zerolog.Logger) *CardDAO {
	logger := log.With().Str("module", "CardDAO").Logger()
	collection := db.Collection("cards")
	return &CardDAO{
		collection: collection,
		log:        logger,
	}
}

func (d *CardDAO) InsertCards(ctx context.Context, cards []models.Card) error {
	logger := d.log.With().Str("method", "insertCard").Logger()
	logger.Info().Msgf("Inserting cards %v", cards)

	c := make([]interface{}, 0, len(cards))
	for _, card := range cards {
		c = append(c, card)
	}

	_, err := d.collection.InsertMany(ctx, c)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting cards %v", cards)
		return errors.Join(fmt.Errorf("error inserting cards: %w", err), ErrInsert)
	}

	return nil
}

func (d *CardDAO) UpdateCard(ctx context.Context, card models.Card) error {
	logger := d.log.With().Str("method", "updateCards").Logger()

	filter := bson.D{{"_id", card.ID}}

	u, err := d.collection.UpdateOne(ctx, filter, bson.M{
		"$set": card,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("Updating card %v", card)
		return errors.Join(fmt.Errorf("error updating card: %w", err), ErrUpdate)
	}

	logger.Info().Msgf("Updated: %+v", u)

	return nil
}

func (d *CardDAO) GetFrontOfCardByID(ctx context.Context, deckID, cardID string) (models.FrontOfCard, error) {
	logger := d.log.With().Str("method", "GetFrontOfCardByID").Logger()
	logger.Info().Msgf("getting front of card by id: %s", cardID)

	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"deck_id", "a415cedb-a732-45cc-a429-1bd3203cd23d"}}}},
		bson.D{
			{"$setWindowFields",
				bson.D{
					{"sortBy", bson.D{{"created_at", -1}}},
					{"output",
						bson.D{
							{"previousCard",
								bson.D{
									{"$push", "$$ROOT"},
									{"window",
										bson.D{
											{"documents",
												bson.A{
													1,
													1,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$setWindowFields",
				bson.D{
					{"sortBy", bson.D{{"created_at", 1}}},
					{"output",
						bson.D{
							{"nextCard",
								bson.D{
									{"$push", "$$ROOT"},
									{"window",
										bson.D{
											{"documents",
												bson.A{
													1,
													1,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$match", bson.D{{"_id", "93dc5755-ce3e-40bc-b974-78f0288419d1"}}}},
		bson.D{
			{"$set",
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
					{"is_upvoted_by_user",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$in",
												bson.A{
													"marken.ryan@gmail.com",
													"$user_upvotes",
												},
											},
										},
									},
									{"then", true},
									{"else", false},
								},
							},
						},
					},
					{"is_downvoted_by_user",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$in",
												bson.A{
													"marken.ryan@gmail.com",
													"$user_downvotes",
												},
											},
										},
									},
									{"then", true},
									{"else", false},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"card_id", "$_id"},
					{"content", "$front"},
					{"deck_id", "$deck_id"},
					{"previous_card", bson.D{{"$first", "$previousCard._id"}}},
					{"next_card", bson.D{{"$first", "$nextCard._id"}}},
					{"upvotes", bson.D{{"$size", "$user_upvotes"}}},
					{"downvotes", bson.D{{"$size", "$user_downvotes"}}},
				},
			},
		},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting cursor")
		return models.FrontOfCard{}, errors.Join(err, ErrAggregate)
	}

	var res []models.FrontOfCard
	err = cursor.All(ctx, &res)
	if err != nil {
		logger.Error().Err(err).Msgf("while unmarshalling to FrontOfCard")
		return models.FrontOfCard{}, errors.Join(err, ErrAggregate)
	}
	if len(res) == 0 {
		return models.FrontOfCard{}, ErrNoResults
	}

	return res[0], nil

}

func (d *CardDAO) GetBackOfCardByID(ctx context.Context, deckID, cardID string) (models.BackOfCard, error) {
	logger := d.log.With().Str("method", "GetBackOfCardByID").Logger()
	logger.Info().Msgf("getting back of card by id: %s", cardID)

	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"deck_id", deckID}}}},
		bson.D{
			{"$setWindowFields",
				bson.D{
					{"sortBy", bson.D{{"created_at", -1}}},
					{"output",
						bson.D{
							{"previousCard",
								bson.D{
									{"$push", "$$ROOT"},
									{"window",
										bson.D{
											{"documents",
												bson.A{
													1,
													1,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$setWindowFields",
				bson.D{
					{"sortBy", bson.D{{"created_at", 1}}},
					{"output",
						bson.D{
							{"nextCard",
								bson.D{
									{"$push", "$$ROOT"},
									{"window",
										bson.D{
											{"documents",
												bson.A{
													1,
													1,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$match", bson.D{{"_id", cardID}}}},
		bson.D{
			{"$set",
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
					{"is_upvoted_by_user",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$in",
												bson.A{
													"marken.ryan@gmail.com",
													"$user_upvotes",
												},
											},
										},
									},
									{"then", true},
									{"else", false},
								},
							},
						},
					},
					{"is_downvoted_by_user",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$in",
												bson.A{
													"marken.ryan@gmail.com",
													"$user_downvotes",
												},
											},
										},
									},
									{"then", true},
									{"else", false},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"card_id", "$_id"},
					{"answer", "$back"},
					{"deck_id", "$deck_id"},
					{"next_card", bson.D{{"$first", "$nextCard._id"}}},
					{"is_upvoted_by_user", "$is_upvoted_by_user"},
					{"is_downvoted_by_user", "$is_downvoted_by_user"},
					{"created_at", "$created_at"},
					{"update_at", "$update_at"},
				},
			},
		},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting cursor")
		return models.BackOfCard{}, errors.Join(err, ErrAggregate)
	}

	var res []models.BackOfCard
	err = cursor.All(ctx, &res)
	if err != nil {
		logger.Error().Err(err).Msgf("while unmarshalling to BackOfCard")
		return models.BackOfCard{}, errors.Join(err, ErrAggregate)
	}
	if len(res) == 0 {
		return models.BackOfCard{}, ErrNoResults
	}

	return res[0], nil
}
