package provider

import (
	"context"
	"github.com/rmarken/reptr/service/internal/database"
	mocks "github.com/rmarken/reptr/service/internal/database/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("should return instance of controller", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepository(ctrl)
		c := New(zerolog.Nop(), repo)

		assert.NotNil(t, c)
	})
}
func TestController_UserExists(t *testing.T) {
	testCases := map[string]struct {
		mockDB     func(mockDB *mocks.MockRepository)
		err        error
		subject    string
		expectID   string
		userExists bool
	}{
		"user exists": {
			mockDB: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("123", nil)
			},
			err:        nil,
			subject:    "exists",
			expectID:   "123",
			userExists: true,
		},
		"user does not exist": {
			mockDB: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("", database.ErrNoResults)
			},
			err:        nil,
			subject:    "notExists",
			expectID:   "",
			userExists: false,
		},
		"error handling": {
			mockDB: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("", mongo.ErrNilCursor)
			},
			err:        mongo.ErrNilCursor,
			subject:    "error",
			expectID:   "",
			userExists: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockRepository(ctrl)
			if tc.mockDB != nil {
				tc.mockDB(mockRepo)
			}

			c := Logic{
				logger: zerolog.Nop(),
				repo:   mockRepo,
			}

			_, exists, err := c.UserExists(context.Background(), tc.subject)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, exists, tc.userExists)
		})
	}
}

func TestController_InsertPair(t *testing.T) {

	// Testing InsertPair
	insertTestCases := map[string]struct {
		mockRepo func(repo *mocks.MockRepository)
		wantErr  error
	}{
		"user does not exist": {
			mockRepo: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("", database.ErrNoResults)
				repo.EXPECT().InsertUserSubjectPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		"user exists": {
			mockRepo: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("123", nil)
			},
			wantErr: ErrUserExists,
		},
		"error handling from first get user": {
			mockRepo: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("", database.ErrFind)
			},
			wantErr: database.ErrFind,
		},
		"error handling from insert  user": {
			mockRepo: func(mockDB *mocks.MockRepository) {
				mockDB.EXPECT().GetUserIDFor(gomock.Any(), gomock.Any()).Return("", database.ErrNoResults)
				mockDB.EXPECT().InsertUserSubjectPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(database.ErrInsert)
			},
			wantErr: database.ErrInsert,
		},
	}

	for name, tc := range insertTestCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockRepository(ctrl)
			if tc.mockRepo != nil {
				tc.mockRepo(mockRepo)
			}

			c := Logic{
				logger: zerolog.Nop(),
				repo:   mockRepo,
			}

			err := c.InsertPair(context.Background(), "123")
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
