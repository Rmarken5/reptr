package deck

import (
	"context"
	"fmt"
	"github.com/rmarken/reptr/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

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
			wantErr: ErrInsert, // Your expected error,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DAO{collection: mt.Coll, log: zerolog.Nop()}
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
					// Add more Card instances if needed
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
				// Set up mock responses here as required
				// For example, you can use mt.Coll.InsertMany to prepopulate data

				docs := make([]bson.D, 0)
				for _, wc := range haveWithCards {
					b, err := bson.Marshal(&wc)
					require.NoError(mt, err)

					var res bson.D
					err = bson.Unmarshal(b, &res)
					require.NoError(mt, err)
					docs = append(docs, res)
				}

				cursor := mtest.CreateCursorResponse(0, "0.0", mtest.FirstBatch, docs...)
				mt.AddMockResponses(cursor)
			},
			wantWithCards: haveWithCards,
			wantErr:       nil,
		},
		//"should return error if aggregation fails": {
		//	from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		//	to:     nil,
		//	limit:  10,
		//	offset: 0,
		//	mockDatabase: func(mt *mtest.T) {
		//		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
		//			Code:    12345, // Simulated error code
		//			Message: "aggregation error",
		//		}))
		//	},
		//	wantWithCards: nil,
		//	wantErr:       ErrAggregate, // Your expected error,
		//},
		//"should handle empty result set": {
		//	from:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		//	to:     nil,
		//	limit:  10,
		//	offset: 0,
		//	mockDatabase: func(mt *mtest.T) {
		//		// Set up mock responses here as required
		//		// For example, you can set an empty cursor
		//	},
		//	wantWithCards: []models.WithCards{},
		//	wantErr:       nil,
		//},
		// Add more test cases to cover other scenarios...
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			dao := DAO{collection: mt.Coll, log: zerolog.Nop()}
			if tc.mockDatabase != nil {
				tc.mockDatabase(mt)
			}

			gotWithCards, gotErr := dao.GetWithCards(context.Background(), tc.from, tc.to, tc.limit, tc.offset)
			fmt.Printf("%+v\n%+v\n", gotWithCards, tc.wantWithCards)
			assert.Equal(t, tc.wantErr, gotErr)
			assert.Len(t, gotWithCards, len(tc.wantWithCards))
			for i, gotWithCard := range gotWithCards {
				assert.Equal(t, tc.wantWithCards[i], gotWithCard)
			}
		})
	}
}
