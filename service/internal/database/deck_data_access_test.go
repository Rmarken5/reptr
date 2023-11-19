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

func TestNewDeckDataAccess(t *testing.T) {
	// Create a mock database and logger for testing
	var (
		mockDB  = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mockLog = zerolog.Nop()
	)
	defer mockDB.Close()

	testLogger := mockLog.With().Str("module", "CardDAO").Logger()

	// Test case to verify NewDataAccess function
	mockDB.Run("NewDataAccess", func(t *mtest.T) {
		dao := NewDeckDataAccess(t.DB, testLogger)

		assert.NotNil(t, dao)
		assert.Equal(t, "decks", dao.collection.Name())
	})
}

func TestDAO_InsertDeck(t *testing.T) {

	var (
		db         = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		writeError = mtest.WriteError{Code: 11000, Message: "duplicate key error"}
	)
	defer db.Close()

	testCases := map[string]struct {
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should return id when document is inserted": {
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		"should return error on database operation failure": {
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(writeError))
			},
			wantErr: ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			_, gotErr := dao.InsertDeck(context.Background(), models.Deck{})
			assert.ErrorIs(mt, gotErr, tc.wantErr)
		})
	}
}

func TestDAO_GetWithCards(t *testing.T) {
	var (
		db            = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		haveWithCards = []models.WithCards{
			{
				Deck: models.Deck{
					ID:        "1", // Populate with appropriate ID
					Name:      "Sample Deck",
					CreatedAt: time.Now().In(time.UTC).Truncate(time.Millisecond),
				},
				Cards: []models.Card{
					{
						ID:        "1", // Populate with appropriate ID
						Front:     "Front Text",
						Back:      "Back Text",
						Kind:      1,   // Set appropriate Type
						DeckID:    "1", // Match the Deck ID above
						CreatedAt: time.Now().In(time.UTC).Truncate(time.Millisecond),
						UpdatedAt: time.Now().In(time.UTC).Truncate(time.Millisecond),
					},
				},
			}}
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		from   time.Time
		to     *time.Time
		limit  int
		offset int

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantWithCards []models.WithCards
		wantErr       error
	}{
		"should get WithCards successfully": {
			from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     nil,
			limit:  10,
			offset: 0,
			mockDatabase: func(mt *mtest.T) {
				docs := make([]bson.D, 0)
				for _, wc := range haveWithCards {
					b, err := bson.Marshal(&wc)
					require.NoError(mt, err)

					var res bson.D
					err = bson.Unmarshal(b, &res)
					require.NoError(mt, err)
					docs = append(docs, res)
				}

				cursor := mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", "dbName", "collName"), mtest.FirstBatch, docs...)
				mt.AddMockResponses(cursor)
			},
			wantWithCards: haveWithCards,
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
			wantWithCards: nil,
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
			wantWithCards: []models.WithCards{},
			wantErr:       ErrNoResults,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotWithCards, gotErr := dao.GetWithCards(context.Background(), tc.from, tc.to, tc.limit, tc.offset)
			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Len(t, gotWithCards, len(tc.wantWithCards))
			for i, gotWithCard := range gotWithCards {
				assert.Equal(t, tc.wantWithCards[i], gotWithCard)
			}
		})
	}
}

func TestDAO_AddUserToUpvote(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		deckID string
		userID string

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should add user to upvote successfully": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if update fails": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
					Message: "update error",
				}))
			},
			wantErr: ErrUpdate,
		},
		// Add more test cases to cover other scenarios...
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.AddUserToUpvote(context.Background(), tc.deckID, tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestDAO_RemoveUserFromUpvote(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		deckID string
		userID string

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should remove user from upvote successfully": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where user is in upvote list and removal succeeds
				// Set up the mock response accordingly
				mt.Coll.InsertOne(context.Background(), bson.M{
					"_id":           "1",
					"user_upvote":   []string{"user123"}, // Existing user in upvote list
					"user_downvote": []string{},
				})
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if update fails": {
			deckID: "1",
			userID: "user123",
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
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
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.RemoveUserFromUpvote(context.Background(), tc.deckID, tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestDAO_AddUserToDownvote(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		deckID string
		userID string

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should add user to downvote successfully": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where user is successfully added to downvote list
				// Set up the mock response accordingly
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if update fails": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where update fails
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
					Message: "update error",
				}))
			},
			wantErr: ErrUpdate, // Or your expected error for this scenario
		},
		// Add more test cases to cover other scenarios...
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.AddUserToDownvote(context.Background(), tc.deckID, tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestDAO_RemoveUserFromDownvote(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		deckID string
		userID string

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should remove user from downvote successfully": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where user is successfully removed from downvote list
				// Set up the mock response accordingly
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if update fails": {
			deckID: "1",       // Set deck ID
			userID: "user123", // Set user ID
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where update fails
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
			dao := DeckDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.RemoveUserFromDownvote(context.Background(), tc.deckID, tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
