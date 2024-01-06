package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	mockAuth "github.com/rmarken/reptr/service/internal/logic/auth/mocks"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	mockLogic "github.com/rmarken/reptr/service/internal/logic/decks/mocks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	r := New(zerolog.Nop(), &decks.Logic{}, &provider.Logic{}, &auth.Authenticator{})
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

func TestReprtClient_Login(t *testing.T) {
	testCases := map[string]struct {
		haveHeader        http.Header
		mockAuthenticator func(mock *mockAuth.MockAuthentication)
		wantStatus        int
		wantResponse      string
	}{
		"should return token on login": {
			haveHeader: http.Header{"Authorization": []string{"Basic am9obi1kb2U6c2VjcmV0cEBzc3cwcmQ="}},
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {
				extra := map[string]interface{}{IDToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkZha2VUZXN0IFVzZXIiLCJleHAiOjE2NDk5OTQ0MDB9.wrongsignature"}
				token := &oauth2.Token{
					AccessToken:  "",
					TokenType:    "",
					RefreshToken: "",
					Expiry:       time.Time{},
				}
				token = token.WithExtra(extra)
				mock.EXPECT().PasswordCredentialsToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(token, nil)
			},
			wantResponse: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkZha2VUZXN0IFVzZXIiLCJleHAiOjE2NDk5OTQ0MDB9.wrongsignature",
			wantStatus:   http.StatusOK,
		},
		"should return token on login with colon in pass": {
			haveHeader: http.Header{"Authorization": []string{"Basic am9obi1kb2U6bXk6cGFzc3dvcmQxMjM="}},
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {
				extra := map[string]interface{}{IDToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkZha2VUZXN0IFVzZXIiLCJleHAiOjE2NDk5OTQ0MDB9.wrongsignature"}
				token := &oauth2.Token{
					AccessToken:  "",
					TokenType:    "",
					RefreshToken: "",
					Expiry:       time.Time{},
				}
				token = token.WithExtra(extra)
				mock.EXPECT().PasswordCredentialsToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(token, nil)
			},
			wantResponse: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkZha2VUZXN0IFVzZXIiLCJleHAiOjE2NDk5OTQ0MDB9.wrongsignature",
			wantStatus:   http.StatusOK,
		},
		"should return unauthorized when authenticator does not return token": {
			haveHeader: http.Header{"Authorization": []string{"Basic am9obi1kb2U6c2VjcmV0cEBzc3cwcmQ="}},
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {

				mock.EXPECT().PasswordCredentialsToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("random error"))
			},
			wantResponse: "",
			wantStatus:   http.StatusUnauthorized,
		},
		"should return bad request when header isn't basic auth": {
			haveHeader:   http.Header{"Authorization": []string{"Bearer am9obi1kb2U6c2VjcmV0cEBzc3cwcmQ="}},
			wantResponse: "",
			wantStatus:   http.StatusBadRequest,
		},
		"should return bad request when header doesn't contain password": {
			haveHeader:   http.Header{"Authorization": []string{"Basic am9obi1kb2U6"}},
			wantResponse: "",
			wantStatus:   http.StatusBadRequest,
		},
		"should return bad request when header doesn't contain user": {
			haveHeader:   http.Header{"Authorization": []string{"Basic OnBhc3N3b3Jk"}},
			wantResponse: "",
			wantStatus:   http.StatusBadRequest,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockAuthenticator := mockAuth.NewMockAuthentication(ctrl)
			if tc.mockAuthenticator != nil {
				tc.mockAuthenticator(mockAuthenticator)
			}
			reprt := ReprtClient{
				logger:        zerolog.Nop(),
				authenticator: mockAuthenticator,
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodGet, "/login", nil)
			require.NoError(t, err)

			req.Header = tc.haveHeader

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.Login(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			if rr.Code == http.StatusOK {
				assert.Equal(t, tc.wantResponse, rr.Body.String())
				require.NoError(t, err)
			}
		})
	}
}
