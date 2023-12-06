package database

import (
	"context"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	var (
		db       = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		logger   = zerolog.Nop()
		userID   = "1" // Set an appropriate user ID
		haveUser = models.User{
			ID:       userID,
			Username: "testuser",
			// Add other necessary user fields
		}
	)
	defer db.Close()

	testCases := map[string]struct {
		userID       string
		mockDatabase func(mongo *mtest.T)
		wantUser     models.User
		wantErr      error
	}{
		"should get user by ID successfully": {
			userID: userID,
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
			userID: userID,
			mockDatabase: func(mt *mtest.T) {
				// Set up mock response for no documents found
				cursor := mtest.CreateCursorResponse(0, "dbName.collName", mtest.FirstBatch)
				mt.AddMockResponses(cursor)
			},
			wantUser: models.User{},
			wantErr:  ErrNoResults,
		},
		"should return error if FindOne operation fails": {
			userID: userID,
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

			gotUser, gotErr := dao.GetUserByID(context.Background(), tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}

func TestInsertUser(t *testing.T) {
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
