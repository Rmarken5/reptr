package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestProviderUsersDAO_GetUserIDFor(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		wantUserID   string
		mockDatabase func(mongo *mtest.T)
		wantErr      error
	}{
		"should get user for subject": {
			mockDatabase: func(mongo *mtest.T) {
				type result struct {
					UserID string `bson:"user_id"`
				}
				r := result{UserID: "123"}

				b, err := bson.Marshal(&r)
				require.NoError(mongo, err)

				var res bson.D
				err = bson.Unmarshal(b, &res)
				require.NoError(mongo, err)

				c := mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", "dbName", "collName"), mtest.FirstBatch, []bson.D{res}...)
				mongo.AddMockResponses(c)
			},
			wantUserID: "123",
		},
		"should return no results": {
			mockDatabase: func(mongo *mtest.T) {
				mongo.AddMockResponses(mtest.CreateCursorResponse(
					0,
					fmt.Sprintf("%s.%s", "dbName", "collName"),
					mtest.FirstBatch,
				))
			},
			wantUserID: "",
			wantErr:    ErrNoResults,
		},
		"should return error from db": {
			mockDatabase: func(mongo *mtest.T) {
				mongo.AddMockResponses(bson.D{
					{Key: "ok", Value: -1},
				})
			},
			wantErr: ErrFind,
		},
	}

	for name, tc := range testCases {
		db.Run(name, func(t *mtest.T) {
			if tc.mockDatabase != nil {
				tc.mockDatabase(t)
			}
			dao := ProviderUsersDAO{
				collection: t.Coll,
				log:        zerolog.Nop(),
			}
			userID, err := dao.GetUserIDFor(context.Background(), "")
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, userID, tc.wantUserID)
		})
	}
}

func TestProviderUsersDAO_InsertUserSubjectPair(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer db.Close()

	testCases := map[string]struct {
		mockDB  func(mt *mtest.T)
		wantErr error
	}{
		"should return no error on successful insert": {
			mockDB: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		"should return error on unsuccessful insert": {
			mockDB: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Code:    12345, // Simulated error code
					Message: "insert error",
				}))
			},
			wantErr: ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		db.Run(name, func(t *mtest.T) {

			if tc.mockDB != nil {
				tc.mockDB(t)
			}
			dao := ProviderUsersDAO{
				collection: t.Coll,
				log:        zerolog.Nop(),
			}
			err := dao.InsertUserSubjectPair(context.Background(), "", "")
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}

}

func TestNewProviderUsersDataAccess(t *testing.T) {
	var (
		db = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)

	defer db.Close()

	db.Run("should return instance of dao", func(t *mtest.T) {
		dao := NewProviderUsersDataAccess(t.DB, zerolog.Nop())
		assert.NotNil(db, dao)
	})
}
