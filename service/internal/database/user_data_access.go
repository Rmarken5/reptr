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

const usersCollection = "users"

var _ UserDataAccess = new(UserDAO)

type (
	UserDataAccess interface {
		InsertUser(ctx context.Context, user models.User) (string, error)
		GetUserByUsername(ctx context.Context, username string) (models.User, error)
		GetGroupsForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.Group, error)
		AddUserAsMemberOfGroup(ctx context.Context, username string, groupName string) error
	}

	UserDAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

func NewUserDataAccess(db *mongo.Database, log zerolog.Logger) *UserDAO {
	logger := log.With().Str("module", "UserDAO").Logger()
	collection := db.Collection(usersCollection)
	return &UserDAO{
		collection: collection,
		log:        logger,
	}
}

func (u *UserDAO) InsertUser(ctx context.Context, user models.User) (string, error) {
	logger := u.log.With().Str("method", "insertUser").Logger()
	logger.Info().Msgf("Inserting user: %v", user)

	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msgf("Inserting user %v", user)
		return "", errors.Join(err, ErrInsert)
	}

	prim, ok := res.InsertedID.(string)
	if !ok {
		logger.Error().Msgf("cannot return object id from inserted user")
		return "", errors.Join(fmt.Errorf("error inserting user: %w", err), ErrInsert)
	}

	logger.Debug().Msgf("%s", prim)

	return prim, nil
}

func (u *UserDAO) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	logger := u.log.With().Str("method", "getUserByID").Logger()

	one := u.collection.FindOne(ctx, bson.D{{"_id", username}})
	if one.Err() != nil {
		logger.Error().Err(one.Err()).Msgf("getting user by haveUsername %s", username)
		if errors.Is(one.Err(), mongo.ErrNoDocuments) {
			return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", one.Err()), ErrNoResults)
		}
		return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", one.Err()), ErrFind)
	}

	var usr models.User
	err := one.Decode(&usr)
	if err != nil {
		logger.Error().Err(err).Msgf("decoding result to user %s", username)
		return models.User{}, errors.Join(fmt.Errorf("error getting user: %w", err))
	}

	return usr, nil
}

func (u *UserDAO) AddUserAsMemberOfGroup(ctx context.Context, username string, groupName string) error {
	logger := u.log.With().Str("method", "AddUserAsMemberOfGroup").Logger()
	logger.Info().Msgf("adding user %s as member of group %s", username, groupName)

	_, err := u.collection.UpdateOne(ctx, bson.D{{"_id", username}}, bson.D{{"$push", bson.D{{"memberOfGroups", groupName}}}})
	if err != nil {
		err = errors.Join(err, ErrUpdate)
		logger.Error().Err(err).Msgf("while adding user to group: %s - %s", username, groupName)
		return err
	}
	return nil
}

func (u *UserDAO) GetGroupsForUser(ctx context.Context, username string, from time.Time, to *time.Time, limit, offset int) ([]models.Group, error) {
	logger := u.log.With().Str("method", "GetGroupsWithDecksByUser").Logger()
	logger.Info().Msgf("getting groups for user: %s", username)

	matchUser := bson.D{{"$match", bson.D{{"_id", username}}}}
	lookupGroups := bson.D{
		{"$lookup",
			bson.D{
				{"from", "groups"},
				{"localField", "member_of_groups"},
				{"foreignField", "name"},
				{"as", "groups"},
			},
		},
	}
	projectGroups := bson.D{
		{"$project",
			bson.D{
				{"_id", 0},
				{"groups", 1},
			},
		},
	}
	unwindGroups := bson.D{
		{"$unwind",
			bson.D{
				{"path", "$groups"},
				{"preserveNullAndEmptyArrays", true},
			},
		},
	}
	flattenGroup := bson.D{
		{"$project",
			bson.D{
				{"name", "$groups.name"},
				{"created_at", "$groups.created_at"},
				{"updated_at", "$groups.updated_at"},
				{"deleted_at", "$groups.deleted_at"},
				{"deck_ids", "$groups.deck_ids"},
			},
		},
	}

	filter := append([]bson.D{matchUser,
		lookupGroups,
		projectGroups,
		unwindGroups,
		flattenGroup,
	},
		pipeline.Paginate(from, to, limit, offset)...)

	logger.Debug().Msgf("%+v", filter)

	cur, err := u.collection.Aggregate(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Msg("error in calling aggregation")
		return nil, errors.Join(err, ErrAggregate)
	}

	var groups []models.Group
	err = cur.All(ctx, &groups)
	if err != nil {
		logger.Error().Err(err).Msg("while unmarshalling into slice")
		return nil, errors.Join(err, ErrAggregate)
	}

	if len(groups) == 0 {
		logger.Info().Msgf("No groups for %s", username)
		return nil, ErrNoResults
	}
	return groups, nil
}
