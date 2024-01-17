package database

import (
	"context"
	"fmt"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestInsertGroup(t *testing.T) {
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		group        models.Group
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should insert group successfully": {
			group: models.Group{},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error on insertion failure": {
			group: models.Group{},
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

			dao := GroupDAO{
				collection: mt.Coll,
				log:        logger,
			}

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}
			_, gotErr := dao.InsertGroup(context.Background(), tc.group)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	// Set up mock database and logger
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		// Input
		group models.Group

		// Mocks
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should update group successfully": {
			group: models.Group{
				// Initialize with appropriate data
			},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse()) // Simulate update success
			},
			wantErr: nil,
		},
		"should return error on update failure": {
			group: models.Group{
				// Initialize with appropriate data
			},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
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
			dao := NewGroupDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.UpdateGroup(context.Background(), tc.group)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestGetGroupsWithDecks(t *testing.T) {
	var (
		db            = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger        = zerolog.Nop()
		haveWithDecks = []models.GroupWithDecks{
			{
				Group: models.Group{
					ID:        "1",
					Name:      "Group 1",
					CreatedAt: time.Now().UTC().AddDate(0, 0, -7).Truncate(time.Millisecond),
					UpdatedAt: time.Now().UTC().AddDate(0, 0, -5).Truncate(time.Millisecond),
				},
				Decks: []models.GetDeckResults{
					{
						ID:        "1",
						Name:      "Deck 1",
						CreatedAt: time.Now().UTC().AddDate(0, 0, -3).Truncate(time.Millisecond),
						UpdatedAt: time.Now().UTC().AddDate(0, 0, -2).Truncate(time.Millisecond),
					},
					{
						ID:        "2",
						Name:      "Deck 2",
						CreatedAt: time.Now().UTC().AddDate(0, 0, -2).Truncate(time.Millisecond),
						UpdatedAt: time.Now().UTC().AddDate(0, 0, -1).Truncate(time.Millisecond),
					},
				},
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
		wantWithDecks []models.GroupWithDecks
		wantErr       error
	}{
		"should get group with decks successfully": {
			from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     nil,
			limit:  10,
			offset: 0,
			mockDatabase: func(mt *mtest.T) {

				docs := make([]bson.D, 0)
				for _, wd := range haveWithDecks {
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
			wantWithDecks: haveWithDecks,
			wantErr:       nil,
		},
		"should return error if aggregation fails": {
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
			wantWithDecks: nil,
			wantErr:       ErrAggregate,
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
			wantWithDecks: []models.GroupWithDecks{},
			wantErr:       ErrNoResults,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := NewGroupDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotWithDecks, gotErr := dao.GetGroupsWithDecks(context.Background(), tc.from, tc.to, tc.limit, tc.offset)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
			assert.Len(mt, gotWithDecks, len(tc.wantWithDecks))
			for i, wd := range gotWithDecks {
				assert.Equal(mt, tc.wantWithDecks[i], wd)
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	// Set up mock database and logger
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		groupID      string
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should delete group successfully": {
			groupID: "1", // Set appropriate group ID
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for successful deletion
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error on deletion failure": {
			groupID: "2", // Set appropriate group ID
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for deletion failure
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
					Message: "delete error",
				}))
			},
			wantErr: ErrDelete,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := NewGroupDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.DeleteGroup(context.Background(), tc.groupID)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestGetGroupByName(t *testing.T) {
	var (
		db           = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger       = zerolog.Nop()
		haveWithDeck = models.GroupWithDecks{
			Group: models.Group{
				ID:        "1",
				Name:      "Group 1",
				CreatedAt: time.Now().UTC().AddDate(0, 0, -7).Truncate(time.Millisecond),
				UpdatedAt: time.Now().UTC().AddDate(0, 0, -5).Truncate(time.Millisecond),
			},
		}
	)
	defer db.Close()

	testCases := map[string]struct {
		groupName    string
		mockDatabase func(mongo *mtest.T)
		wantGroup    models.GroupWithDecks
		wantErr      error
	}{
		"should get group by name successfully": {
			groupName: "Test Group", // Set appropriate group name
			mockDatabase: func(mt *mtest.T) {

				docs := make([]bson.D, 0)
				b, err := bson.Marshal(&haveWithDeck)
				require.NoError(mt, err)

				var res bson.D
				err = bson.Unmarshal(b, &res)
				require.NoError(mt, err)
				docs = append(docs, res)

				cursor := mtest.CreateCursorResponse(0, "dbName.collName", mtest.FirstBatch, docs...)
				mt.AddMockResponses(cursor)
			},
			wantGroup: haveWithDeck,
			wantErr:   nil,
		},
		"should get group by name without records": {
			groupName: "Test Group", // Set appropriate group name
			mockDatabase: func(mt *mtest.T) {

				docs := make([]bson.D, 0)
				cursor := mtest.CreateCursorResponse(0, "dbName.collName", mtest.FirstBatch, docs...)
				mt.AddMockResponses(cursor)
			},
			wantGroup: models.GroupWithDecks{},
			wantErr:   ErrNoResults,
		},
		"should return error if aggregation fails": {
			groupName: "Nonexistent Group", // Set appropriate nonexistent group name
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for aggregation failure
				mt.AddMockResponses(
					bson.D{
						{Key: "ok", Value: -1},
					},
				)
			},
			wantGroup: models.GroupWithDecks{},
			wantErr:   ErrAggregate,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := NewGroupDataAccess(mt.DB, logger)

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotGroup, gotErr := dao.GetGroupByID(context.Background(), tc.groupName)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
			assert.Equal(mt, tc.wantGroup, gotGroup)
		})
	}
}

func TestAddDeckToGroup(t *testing.T) {
	// Set up mock database and logger
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		groupID string
		deckID  string

		// Mocks
		mockDatabase func(mongo *mtest.T)

		// Expected Result
		wantErr error
	}{
		"should add deck to group successfully": {
			groupID: "group123", // Set appropriate group ID
			deckID:  "deck456",  // Set appropriate deck ID
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for successful addition
				mt.Coll.InsertOne(context.Background(), bson.M{
					"_id":      "deck456", // Match the deck ID above
					"deck_ids": []string{"group123"},
				})
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error on update failure": {
			groupID: "group789", // Set appropriate group ID
			deckID:  "deck101",  // Set appropriate deck ID
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for update failure
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
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
			dao := GroupDAO{collection: mt.Coll, log: logger}

			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.AddDeckToGroup(context.Background(), tc.groupID, tc.deckID)

			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}
