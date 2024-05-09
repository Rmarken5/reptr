package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestUserDAO_GetUserByUsername(t *testing.T) {
	var (
		db           = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger       = zerolog.Nop()
		haveUsername = "testuser"
		haveUser     = models.User{
			Username: haveUsername,
		}
	)
	defer db.Close()

	testCases := map[string]struct {
		haveUsername string
		mockDatabase func(mongo *mtest.T)
		wantUser     models.User
		wantErr      error
	}{
		"should get user by haveUsername successfully": {
			haveUsername: haveUsername,
			mockDatabase: func(mt *mtest.T) {
				b, err := bson.Marshal(&haveUser)
				require.NoError(t, err)

				var res bson.D
				err = bson.Unmarshal(b, &res)
				require.NoError(t, err)

				// Set up mock response for successful FindOne operation
				cursor := mtest.CreateCursorResponse(0, "dbName.collName", mtest.FirstBatch, res)
				mt.AddMockResponses(cursor)
			},
			wantUser: haveUser,
			wantErr:  nil,
		},
		"should return error if user not found": {
			haveUsername: haveUsername,
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for no documents found
				cursor := mtest.CreateCursorResponse(0, "dbName.collName", mtest.FirstBatch)
				mt.AddMockResponses(cursor)
			},
			wantUser: models.User{},
			wantErr:  ErrNoResults,
		},
		"should return error if FindOne operation fails": {
			haveUsername: haveUsername,
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for FindOne operation failure
				mt.AddMockResponses(
					bson.D{
						{Key: "ok", Value: -1},
					},
				)
			},
			wantUser: models.User{},
			wantErr:  ErrFind,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := NewUserDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotUser, gotErr := dao.GetUserByUsername(context.Background(), tc.haveUsername)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}

func TestUserDAO_InsertUser(t *testing.T) {
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		user         models.User
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should insert user successfully": {
			user: models.User{},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error on insertion failure": {
			user: models.User{},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    12345,
					Message: "insertion error",
				}))
			},
			wantErr: ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {

			dao := UserDAO{mt.Coll, logger}

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}
			_, gotErr := dao.InsertUser(context.Background(), tc.user)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestUserDAO_AddUserAsMemberOfGroup(t *testing.T) {
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should push id on memberOfGroups": {
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error on update failure": {
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    12345,
					Message: "update error",
				}))
			},
			wantErr: ErrUpdate,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {

			dao := UserDAO{mt.Coll, logger}

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}
			gotErr := dao.AddUserAsMemberOfGroup(context.Background(), uuid.NewString(), uuid.NewString())

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestUserDAO_GetGroupsForUser(t *testing.T) {
	var (
		timeNow    = time.Now().UTC().Truncate(time.Millisecond)
		db         = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger     = zerolog.Nop()
		haveGroups = []models.HomePageGroup{
			{
				ID:         uuid.NewString(),
				Name:       "test-1",
				CreatedBy:  "user1",
				Moderators: []string{"user1"},
				DeckIDs:    []string{},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
			{
				ID:         uuid.NewString(),
				Name:       "test-2",
				CreatedBy:  "user1",
				Moderators: []string{"user1"},
				DeckIDs:    []string{},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
			{
				ID:         uuid.NewString(),
				Name:       "test-3",
				CreatedBy:  "user1",
				Moderators: []string{"user1"},
				DeckIDs:    []string{},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
		}
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		from   time.Time
		to     *time.Time
		limit  int
		offset int

		// Mocks
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantGroups []models.HomePageGroup
		wantErr    error
	}{
		"should get groups successfully": {
			from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     nil,
			limit:  10,
			offset: 0,
			mockDatabase: func(mt *mtest.T) {
				docs := make([]bson.D, 0)
				for _, wd := range haveGroups {
					b, err := bson.Marshal(&wd)
					require.NoError(mt, err)

					var res bson.D
					err = bson.Unmarshal(b, &res)
					require.NoError(mt, err)
					docs = append(docs, res)
				}

				cursor := mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", "dbName", "collName"), mtest.FirstBatch, docs...)
				mt.AddMockResponses(cursor)
			},
			wantGroups: haveGroups,
			wantErr:    nil,
		},
		"should return aggregation error if aggregation fails": {
			from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     nil,
			limit:  10,
			offset: 0,
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(
					bson.D{
						{Key: "ok", Value: -1},
					},
				)
			},
			wantGroups: nil,
			wantErr:    ErrAggregate,
		},
		"should handle empty result set": {
			from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     nil,
			limit:  10,
			offset: 0,
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(
					0,
					fmt.Sprintf("%s.%s", "dbName", "collName"),
					mtest.FirstBatch,
				))
			},
			wantGroups: []models.HomePageGroup(nil),
			wantErr:    ErrNoResults,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := NewUserDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotGroups, gotErr := dao.GetGroupsForUser(context.Background(), uuid.NewString(), tc.from, tc.to, tc.limit, tc.offset)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
			require.Len(mt, gotGroups, len(tc.wantGroups))
			for i, wd := range gotGroups {
				assert.Equal(mt, tc.wantGroups[i], wd)
			}
		})
	}
}
