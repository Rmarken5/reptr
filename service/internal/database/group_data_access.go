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

var _ GroupDataAccess = &GroupDAO{}

type (
	GroupDataAccess interface {
		InsertGroup(ctx context.Context, group models.Group) error
		UpdateGroup(ctx context.Context, group models.Group) error
		GetGroupsWithDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.GroupWithDecks, error)
		DeleteGroup(ctx context.Context, groupID string) error
		AddDeckToGroup(ctx context.Context, groupID, deckID string) error
		GetGroupByName(ctx context.Context, groupName string) (models.GroupWithDecks, error)
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

func (g *GroupDAO) InsertGroup(ctx context.Context, group models.Group) error {
	logger := g.log.With().Str("method", "insertGroup").Logger()
	logger.Info().Msgf("Inserting group %v", group)

	_, err := g.collection.InsertOne(ctx, group)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting group %v", group)
		return errors.Join(fmt.Errorf("error inserting group: %w", err), ErrInsert)
	}

	return nil
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
	logger := g.log.With().Str("method", "GetWithCards").Logger()
	logger.Info().Msgf("Getting GetGroupsWithDecks %v - %v, limit: %d offset %d", from, to, limit, offset)

	lookupCards :=
		bson.D{{"$lookup",
			bson.D{
				{"from", "decks"},
				{"localField", ""},
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

func (g *GroupDAO) DeleteGroup(ctx context.Context, groupID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDAO) AddDeckToGroup(ctx context.Context, groupID, deckID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDAO) GetGroupByName(ctx context.Context, groupName string) (models.GroupWithDecks, error) {
	//TODO implement me
	panic("implement me")
}
