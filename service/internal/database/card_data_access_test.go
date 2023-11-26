package database

import (
	"context"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestNewCardDataAccess(t *testing.T) {
	// Create a mock database and logger for testing
	var (
		mockDB  = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mockLog = zerolog.Nop()
	)
	defer mockDB.Close()

	testLogger := mockLog.With().Str("module", "CardDAO").Logger()

	// Test case to verify NewDataAccess function
	mockDB.Run("NewCardDataAccess", func(t *mtest.T) {
		dao := NewCardDataAccess(t.DB, testLogger)

		assert.NotNil(t, dao)
		assert.Equal(t, "cards", dao.collection.Name())
	})
}

func TestDAO_InsertCards(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		// Inputs
		cards []models.Card

		// Mocks or setup for MongoDB
		mockDatabase func(mongo *mtest.T)

		// Expected Results
		wantErr error
	}{
		"should insert cards successfully": {
			cards: []models.Card{
				{
					ID:        "1", // Set card ID
					Front:     "Front Text",
					Back:      "Back Text",
					Kind:      1,   // Set appropriate Type
					DeckID:    "1", // Match the Deck ID above
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockDatabase: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if insertion fails": {
			cards: []models.Card{
				{},
			},
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where insertion fails
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
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
			dao := CardDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.InsertCards(context.Background(), tc.cards)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestDAO_UpdateCard(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		card         models.Card
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should update card successfully": {
			card: models.Card{
				ID:        "1", // Set card ID
				Front:     "Updated Front Text",
				Back:      "Updated Back Text",
				Kind:      1,   // Set appropriate Type
				DeckID:    "1", // Match the Deck ID above
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where the card is successfully updated
				// Set up the mock response accordingly
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		"should return error if update fails": {
			card: models.Card{
				// Set card data
			},
			mockDatabase: func(mt *mtest.T) {
				// Simulate scenario where the update fails
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
			dao := CardDAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotErr := dao.UpdateCard(context.Background(), tc.card)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
