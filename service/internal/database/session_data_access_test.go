package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestSessionDAO_SetAnswerForCard(t *testing.T) {
	var (
		db     = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger = zerolog.Nop()
	)
	defer db.Close()

	testCases := map[string]struct {
		haveSessionID          string
		haveCardID             string
		hasIsAnsweredCorrectly bool
		mockMongo              func(mongo *mtest.T)
		wantErr                error
	}{
		"add new card answer to collection when card_id does not exist in session array": {
			haveSessionID:          uuid.NewString(),
			haveCardID:             uuid.NewString(),
			hasIsAnsweredCorrectly: true,
			mockMongo: func(mt *mtest.T) {
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(
						bson.E{Key: "n", Value: 0},
						bson.E{Key: "nModified", Value: 0},
					),
					mtest.CreateSuccessResponse(
						bson.E{Key: "n", Value: 1},
						bson.E{Key: "nModified", Value: 1},
					),
				)
			},
		},
		"update card answer in array when card_id exist in session array": {
			haveSessionID:          uuid.NewString(),
			haveCardID:             uuid.NewString(),
			hasIsAnsweredCorrectly: true,
			mockMongo: func(mt *mtest.T) {
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(
						bson.E{Key: "n", Value: 1},
						bson.E{Key: "nModified", Value: 1},
					),
				)
			},
		},
		"should not update answers when session does not exist": {
			haveSessionID:          uuid.NewString(),
			haveCardID:             uuid.NewString(),
			hasIsAnsweredCorrectly: true,
			wantErr:                ErrFind,
			mockMongo: func(mt *mtest.T) {
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(
						bson.E{Key: "n", Value: 0},
						bson.E{Key: "nModified", Value: 0},
					),
					mtest.CreateSuccessResponse(
						bson.E{Key: "n", Value: 0},
						bson.E{Key: "nModified", Value: 0},
					),
				)
			},
		},
		"return error from mongo on first query": {
			haveSessionID:          uuid.NewString(),
			haveCardID:             uuid.NewString(),
			hasIsAnsweredCorrectly: true,
			mockMongo: func(mt *mtest.T) {

				mt.AddMockResponses(mtest.CreateSuccessResponse(
					bson.E{Key: "n", Value: 0},
					bson.E{Key: "nModified", Value: 0},
				),
					mtest.CreateCommandErrorResponse(mtest.CommandError{
						Code:    12345,
						Message: "insertion error",
					}))
			},
			wantErr: ErrUpdate,
		},
		"return error from mongo on second query": {
			haveSessionID:          uuid.NewString(),
			haveCardID:             uuid.NewString(),
			hasIsAnsweredCorrectly: true,
			mockMongo: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    12345,
					Message: "insertion error",
				}))
			},
			wantErr: ErrUpdate,
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(mt *mtest.T) {
			ctx := context.Background()

			if tc.mockMongo != nil {
				tc.mockMongo(mt)
			}

			sessionDAO := SessionDAO{
				collection: mt.Coll,
				log:        logger,
			}

			err := sessionDAO.SetAnswerForCard(ctx, tc.haveSessionID, tc.haveCardID, tc.hasIsAnsweredCorrectly)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
