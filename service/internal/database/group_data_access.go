package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/service/internal/database/pipeline"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	_ GroupDataAccess = &GroupDAO{}

	deckFromGroupsLookup = bson.D{
		{"$lookup",
			bson.D{
				{"from", "decks"},
				{"localField", "deck_ids"},
				{"foreignField", "_id"},
				{"as", "decks"},
			},
		},
	}
)

type (
	GroupDataAccess interface {
		InsertGroup(ctx context.Context, group models.Group) (string, error)
		UpdateGroup(ctx context.Context, group models.Group) error
		GetGroupsWithDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error)
		DeleteGroup(ctx context.Context, groupID string) error
		GetGroupByID(ctx context.Context, groupID string) (models.GroupWithDecks, error)
		AddDeckToGroup(ctx context.Context, groupID, deckID string) error
		// AddUserToGroup(ctx context.Context, groupID string, haveUsername string) error
	}
	GroupDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewGroupDataAccess(db *mongo.Database, log zerolog.Logger) *GroupDAO {
	logger := log.With().Str("module", "GroupDAO").Logger()
	collection := db.Collection("groups")
	return &GroupDAO{
		collection: collection,
		log:        logger,
	}
}

func (g *GroupDAO) InsertGroup(ctx context.Context, group models.Group) (string, error) {
	logger := g.log.With().Str("method", "insertGroup").Logger()
	logger.Info().Msgf("Inserting group %v", group)

	res, err := g.collection.InsertOne(ctx, group)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting group %v", group)
		return "", errors.Join(err, ErrInsert)
	}

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted group")
		return "", errors.Join(fmt.Errorf("error inserting group: %w", err), ErrInsert)
	}

	return prim, nil
}

func (g *GroupDAO) UpdateGroup(ctx context.Context, group models.Group) error {
	logger := g.log.With().Str("method", "updateGroup").Logger()

	filter := bson.D{{"_id", group.ID}}

	u, err := g.collection.UpdateOne(ctx, filter, bson.M{
		"$set": group,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("Updating group %v", group)
		return errors.Join(fmt.Errorf("error updating group: %w", err), ErrUpdate)
	}

	logger.Info().Msgf("Updated: %+v", u)

	return nil
}

func (g *GroupDAO) GetGroupsWithDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error) {
	logger := g.log.With().Str("method", "GetGroupsWithDecks").Logger()
	logger.Info().Msgf("Getting GetGroupsWithDecks %v - %v, limit: %d offset %d", from, to, limit, offset)

	filter := append(
		pipeline.Paginate(from, to, limit, offset),
		deckFromGroupsLookup,
	)

	cur, err := g.collection.Aggregate(ctx, filter)
	if err != nil {
		return nil, errors.Join(err, ErrAggregate)
	}
	defer cur.Close(ctx)

	withDecks := make([]models.GroupWithDecks, 0)
	err = cur.All(ctx, &withDecks)
	if err != nil {
		return nil, err
	}
	if len(withDecks) == 0 {
		return nil, ErrNoResults
	}
	return withDecks, nil
}

func (g *GroupDAO) DeleteGroup(ctx context.Context, groupID string) error {
	logger := g.log.With().Str("method", "deleteGroup").Logger()

	filter := bson.D{{"_id", groupID}}

	u, err := g.collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.D{
			{"deleted_at", time.Now().String()},
		},
	})
	if err != nil {
		logger.Error().Err(err).Msgf("deleting group %s", groupID)
		return errors.Join(fmt.Errorf("error deleting group: %w", err), ErrDelete)
	}

	logger.Info().Msgf("Updated: %+v", u)

	return nil
}

func (g *GroupDAO) AddDeckToGroup(ctx context.Context, groupID, deckID string) error {
	logger := g.log.With().Str("method", "AddDeckToGroup").Logger()
	logger.Info().Msgf("adding group to deck: %s", groupID)

	filter := bson.D{{Key: "_id", Value: groupID}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"deck_ids", deckID},
		}},
	}

	_, err := g.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(fmt.Errorf("adding group to deck: %w", err), ErrUpdate)
	}

	return nil
}

func (g *GroupDAO) GetGroupByID(ctx context.Context, groupID string) (models.GroupWithDecks, error) {
	logger := g.log.With().Str("method", "GetGroupByID").Logger()
	match := bson.D{{"$match", bson.D{{"_id", groupID}}}}

	lookupDecks := bson.D{
		{"$lookup",
			bson.D{
				{"from", "decks"},
				{"localField", "deck_ids"},
				{"foreignField", "_id"},
				{"as", "decks"},
			},
		},
	}
	getVotes := bson.D{
		{"$addFields",
			bson.D{
				{"decks",
					bson.D{
						{"$map",
							bson.D{
								{"input", "$decks"},
								{"as", "deck"},
								{"in",
									bson.D{
										{"$mergeObjects",
											bson.A{
												"$$deck",
												bson.D{
													{"upvotes",
														bson.D{
															{"$size",
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$$deck.user_upvotes",
																			bson.A{},
																		},
																	},
																},
															},
														},
													},
													{"downvotes",
														bson.D{
															{"$size",
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$$deck.user_downvotes",
																			bson.A{},
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
							},
						},
					},
				},
			},
		},
	}
	unwind := bson.D{
		{"$unwind",
			bson.D{
				{"path", "$decks"},
				{"preserveNullAndEmptyArrays", true},
			},
		},
	}
	lookupCards := bson.D{
		{"$lookup",
			bson.D{
				{"from", "cards"},
				{"localField", "decks._id"},
				{"foreignField", "deck_id"},
				{"as", "cards"},
			},
		},
	}

	getSizeOfCards := bson.D{
		{"$addFields",
			bson.D{
				{"decks",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$ne",
											bson.A{
												bson.D{{"$size", "$deck_ids"}},
												0,
											},
										},
									},
								},
								{"then",
									bson.D{
										{"$mergeObjects",
											bson.A{
												"$decks",
												bson.D{
													{"num_cards",
														bson.D{
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
														},
													},
												},
											},
										},
									},
								},
								{"else", primitive.Null{}},
							},
						},
					},
				},
			},
		},
	}

	regroup := bson.D{
		{"$group",
			bson.D{
				{"_id", "$_id"},
				{"name", bson.D{{"$first", "$name"}}},
				{"created_by", bson.D{{"$first", "$created_by"}}},
				{"created_at", bson.D{{"$first", "$created_at"}}},
				{"updated_at", bson.D{{"$first", "$updated_at"}}},
				{"deleted_at", bson.D{{"$first", "$deleted_at"}}},
				{"members", bson.D{{"$first", "$members"}}},
				{"moderators", bson.D{{"$first", "$moderators"}}},
				{"decks",
					bson.D{
						{"$push",
							bson.D{
								{"$cond",
									bson.D{
										{"if",
											bson.D{
												{"$ne",
													bson.A{
														"$decks",
														primitive.Null{},
													},
												},
											},
										},
										{"then", "$decks"},
										{"else", "$$REMOVE"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	removeUserVotes := bson.D{
		{"$project",
			bson.D{
				{"decks.user_upvotes", 0},
				{"decks.user_downvotes", 0},
			},
		},
	}

	filter := bson.A{
		match,
		lookupDecks,
		getVotes,
		unwind,
		lookupCards,
		getSizeOfCards,
		regroup,
		removeUserVotes,
	}

	cur, err := g.collection.Aggregate(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Msgf("getting group by id %s", groupID)
		return models.GroupWithDecks{}, errors.Join(fmt.Errorf("error deleting group: %w", err), ErrAggregate)
	}
	defer cur.Close(ctx)

	withDecks := make([]models.GroupWithDecks, 0)
	err = cur.All(ctx, &withDecks)
	if err != nil {
		return models.GroupWithDecks{}, err
	}
	if len(withDecks) == 0 {
		return models.GroupWithDecks{}, ErrNoResults
	}
	return withDecks[0], nil
}
