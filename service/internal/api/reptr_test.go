package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	mockLogic "github.com/rmarken/reptr/service/internal/logic/decks/mocks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	r := New(zerolog.Nop(), &decks.Logic{})
	assert.NotNil(t, r)
}
func TestGetGroups(t *testing.T) {
	var (
		now                 = time.Now().In(time.UTC)
		nowMinusHour        = now.Add(-1 * time.Hour)
		groupName           = uuid.NewString()
		groupID             = uuid.NewString()
		deckName            = uuid.NewString()
		deckID              = uuid.NewString()
		haveGroupsWithDecks = []models.GroupWithDecks{
			{Group: models.Group{
				ID:        groupID,
				Name:      groupName,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			}, Decks: []models.Deck{
				{ID: deckID, Name: deckName, CreatedAt: now, UpdatedAt: now},
			}},
		}
		wantDecks = []api.Deck{
			{
				CreatedAt: now,
				Id:        deckID,
				Name:      deckName,
				UpdatedAt: now,
			},
		}
		wantGroups = api.GetGroups{
			{
				CreatedAt: now,
				Decks:     wantDecks,
				Id:        groupID,
				Name:      groupName,
				UpdatedAt: now,
			},
		}
	)
	testCases := map[string]struct {
		Params       api.GetGroupsParams
		mockCtrl     func(mock *mockLogic.MockController)
		wantGroups   api.GetGroups
		ExpectedCode int
	}{
		"ValidRequest": {
			Params: api.GetGroupsParams{
				From:   time.Now(),
				Limit:  10,
				Offset: 0,
			},
			ExpectedCode: http.StatusOK,
			mockCtrl: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroups(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(haveGroupsWithDecks, nil)
			},
			wantGroups: wantGroups,
		},
		"Invalid time window": {
			Params: api.GetGroupsParams{
				From:   now,
				To:     &nowMinusHour,
				Limit:  10,
				Offset: 0,
			},
			ExpectedCode: http.StatusBadRequest,
			mockCtrl: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroups(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, decks.ErrInvalidToBeforeFrom)
			},
			wantGroups: nil,
		},
	}

	// Iterate through test cases and perform tests
	for testName, testCase := range testCases {
		testName := testName
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			l := mockLogic.NewMockController(ctrl)

			if testCase.mockCtrl != nil {
				testCase.mockCtrl(l)
			}

			rc := ReprtClient{
				logger:         zerolog.Nop(),
				deckController: l,
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest("GET", "/groups", nil)
			require.NoError(t, err)

			// Set parameters in query string
			q := req.URL.Query()
			q.Add("from", testCase.Params.From.Format(time.RFC3339))
			q.Add("limit", fmt.Sprintf("%d", testCase.Params.Limit))
			q.Add("offset", fmt.Sprintf("%d", testCase.Params.Offset))
			req.URL.RawQuery = q.Encode()

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the GetGroups method with the mocked dependencies
			rc.GetGroups(rr, req, testCase.Params)

			// Check the status code of the response
			assert.Equal(t, testCase.ExpectedCode, rr.Code)
			if rr.Code == http.StatusOK {
				resBytes := rr.Body.Bytes()
				var gotGroups api.GetGroups
				err = json.Unmarshal(resBytes, &gotGroups)
				require.NoError(t, err)
				assert.ElementsMatch(t, testCase.wantGroups, gotGroups)
			}
		})
	}
}
