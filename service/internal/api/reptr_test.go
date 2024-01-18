package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/api"
	mocks "github.com/rmarken/reptr/service/internal/api/middlewares/mocks"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/database"
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
	"net/url"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	r := New(zerolog.Nop(), &decks.Logic{}, &provider.Logic{}, &auth.Authenticator{}, &sessions.CookieStore{})
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
			}, Decks: []models.GetDeckResults{
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

func TestReprtClient_RegistrationPage(t *testing.T) {
	t.Run("return registration page", func(t *testing.T) {
		reprt := ReprtClient{
			logger: zerolog.Nop(),
		}
		// Create a request object with necessary parameters
		req, err := http.NewRequest(http.MethodGet, "/register", nil)
		require.NoError(t, err)

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the Login method with the mocked dependencies
		reprt.RegistrationPage(rr, req)

		// Check the status code of the response
		assert.Equal(t, http.StatusOK, rr.Code)
		snaps.MatchSnapshot(t, rr.Body.String())
	})
}

func TestReprtClient_Register(t *testing.T) {
	testCases := map[string]struct {
		haveEmail         string
		havePassword      string
		haveRepass        string
		mockAuthenticator func(mock *mockAuth.MockAuthentication)
		wantStatus        int
	}{
		"should return return login page on successful registration": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
			haveRepass:   "str0ngP@ssword!",
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {
				mock.EXPECT().
					RegisterUser(gomock.Any(), "someone@somewhere.com", "str0ngP@ssword!").
					Return(models.RegistrationUser{
						ID:            uuid.NewString(),
						Email:         "fake@email.com",
						EmailVerified: false,
					}, models.RegistrationError{}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		"should return bad request when user doesn't provide email": {
			haveEmail:  "",
			wantStatus: http.StatusBadRequest,
		},
		"should return bad request when user doesn't provide password": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "",
			wantStatus:   http.StatusBadRequest,
		},
		"should return bad request when passwords don't match": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
			haveRepass:   "str0ngP@ssword!1",
			wantStatus:   http.StatusBadRequest,
		},
		"should return bad request when validator returns registration error": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
			haveRepass:   "str0ngP@ssword!",
			wantStatus:   http.StatusBadRequest,
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {
				mock.EXPECT().
					RegisterUser(gomock.Any(), "someone@somewhere.com", "str0ngP@ssword!").
					Return(models.RegistrationUser{}, models.RegistrationError{
						Name:        "invalid_registration",
						Code:        "invalid_registration",
						Description: "invalid_registration",
						StatusCode:  http.StatusBadRequest,
					}, nil)
			},
		},
		"should return bad request when validator returns error": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
			haveRepass:   "str0ngP@ssword!",
			wantStatus:   http.StatusInternalServerError,
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {
				mock.EXPECT().
					RegisterUser(gomock.Any(), "someone@somewhere.com", "str0ngP@ssword!").
					Return(models.RegistrationUser{}, models.RegistrationError{}, errors.New("oops"))
			},
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
			req, err := http.NewRequest(http.MethodPost, "/login", nil)
			require.NoError(t, err)

			req.PostForm = map[string][]string{}
			req.PostForm.Set("email", tc.haveEmail)
			req.PostForm.Set("password", tc.havePassword)
			req.PostForm.Set("repassword", tc.haveRepass)

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.Register(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			snaps.MatchSnapshot(t, rr.Body.String())

		})
	}
}

func TestReprtClient_Login(t *testing.T) {
	testCases := map[string]struct {
		haveEmail         string
		havePassword      string
		mockAuthenticator func(mock *mockAuth.MockAuthentication)
		mockStore         func(mock *mocks.MockStore)
		wantStatus        int
		wantResponse      string
	}{
		"should return token on login": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
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
				mock.EXPECT().VerifyIDToken(gomock.Any(), token.Extra(IDToken)).Return(&oidc.IDToken{
					Issuer:          "",
					Audience:        nil,
					Subject:         "123-456",
					Expiry:          time.Time{},
					IssuedAt:        time.Time{},
					Nonce:           "",
					AccessTokenHash: "",
				}, nil)
			},
			mockStore: func(mock *mocks.MockStore) {
				session := sessions.NewSession(mock, uuid.NewString())
				mock.EXPECT().Get(gomock.Any(), CookieSessionID).Return(session, nil)
				mock.EXPECT().Save(gomock.Any(), gomock.Any(), session).Return(nil)
			},
			wantResponse: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkZha2VUZXN0IFVzZXIiLCJleHAiOjE2NDk5OTQ0MDB9.wrongsignature",
			wantStatus:   http.StatusSeeOther,
		},
		"should return unauthorized when authenticator does not return token": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "str0ngP@ssword!",
			mockAuthenticator: func(mock *mockAuth.MockAuthentication) {

				mock.EXPECT().PasswordCredentialsToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("random error"))
			},
			wantResponse: "",
			wantStatus:   http.StatusUnauthorized,
		},
		"should return bad request when header doesn't contain password": {
			haveEmail:    "someone@somewhere.com",
			havePassword: "",
			wantResponse: "",
			wantStatus:   http.StatusBadRequest,
		},
		"should return bad request when header doesn't contain email": {
			haveEmail:    "",
			havePassword: "str0ngP@ssword!",
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
			mockStore := mocks.NewMockStore(ctrl)
			if tc.mockAuthenticator != nil {
				tc.mockAuthenticator(mockAuthenticator)
			}

			if tc.mockStore != nil {
				tc.mockStore(mockStore)
			}

			reprt := ReprtClient{
				logger:        zerolog.Nop(),
				authenticator: mockAuthenticator,
				store:         mockStore,
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodPost, "/login", nil)
			require.NoError(t, err)

			req.PostForm = map[string][]string{}
			req.PostForm.Set("email", tc.haveEmail)
			req.PostForm.Set("password", tc.havePassword)

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.Login(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			gotToken := rr.Header().Get("Authorization")
			assert.Equal(t, tc.wantResponse, gotToken)
			if rr.Code == http.StatusOK {
				snaps.MatchSnapshot(t, rr.Body.String())
				require.NoError(t, err)
			}
		})
	}
}

func TestReprtClient_LoginPage(t *testing.T) {
	t.Run("return login page", func(t *testing.T) {
		reprt := ReprtClient{
			logger: zerolog.Nop(),
		}
		// Create a request object with necessary parameters
		req, err := http.NewRequest(http.MethodGet, "/login", nil)
		require.NoError(t, err)

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the Login method with the mocked dependencies
		reprt.LoginPage(rr, req)

		// Check the status code of the response
		assert.Equal(t, http.StatusOK, rr.Code)
		snaps.MatchSnapshot(t, rr.Body.String())
	})
}

func TestReprtClient_GroupPage(t *testing.T) {
	var (
		timeNow   = time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
		haveGroup = models.GroupWithDecks{
			Group: models.Group{
				ID:        "1234",
				Name:      "name",
				CreatedBy: "someone@somewher.com",
				Moderators: []string{
					"someone@somewher.com",
				},
				DeckIDs: []string{
					"456",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
				DeletedAt: nil,
			},
			Decks: []models.GetDeckResults{
				{
					ID:        "deckID",
					Name:      "deckName",
					Upvotes:   1,
					Downvotes: 3,
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
			},
		}
	)
	testCases := map[string]struct {
		mockController func(mock *mockLogic.MockController)
		wantGroups     models.GroupWithDecks
		wantStatus     int
	}{
		"should load group page with group data": {
			mockController: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).Return(haveGroup, nil)
			},
			wantGroups: haveGroup,
			wantStatus: http.StatusOK,
		},
		"should return 404 when error from database returns not found": {
			mockController: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).Return(models.GroupWithDecks{}, database.ErrNoResults)
			},
			wantGroups: models.GroupWithDecks{},
			wantStatus: http.StatusNotFound,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mockLogic.NewMockController(ctrl)
			if tc.mockController != nil {
				tc.mockController(mock)
			}

			reprt := ReprtClient{
				deckController: mock,
				logger:         zerolog.Nop(),
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodGet, "/page/group/{groupID}", nil)
			require.NoError(t, err)

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.GroupPage(rr, req, uuid.NewString())

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			snaps.MatchSnapshot(t, rr.Body.String())
		})
	}
}

func TestReprtClient_HomePage(t *testing.T) {
	var (
		timeNow   = time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
		haveGroup = []models.Group{
			{
				ID:        "1234",
				Name:      "name",
				CreatedBy: "someone@somewher.com",
				Moderators: []string{
					"someone@somewher.com",
				},
				DeckIDs: []string{
					"456",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
				DeletedAt: nil,
			},
			{
				ID:        "5678",
				Name:      "name-1",
				CreatedBy: "someone@somewhere.com",
				Moderators: []string{
					"someone@somewhere.com",
				},
				DeckIDs: []string{
					"456",
					"789",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
				DeletedAt: nil,
			},
		}
	)
	testCases := map[string]struct {
		mockController func(mock *mockLogic.MockController)
		wantGroups     models.Group
		wantStatus     int
		wantUserName   string
	}{
		"should load group page with group data": {
			mockController: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroupsForUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(haveGroup, nil)
			},
			wantUserName: "hello",
			wantStatus:   http.StatusOK,
		},
		"should return 404 when error from database returns not found": {
			mockController: func(mock *mockLogic.MockController) {
				mock.EXPECT().GetGroupsForUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, database.ErrNoResults)
			},
			wantUserName: "hello",
			wantGroups:   models.Group{},
			wantStatus:   http.StatusNotFound,
		},
		"should return internal error when username is not on context": {
			wantUserName: "",
			wantGroups:   models.Group{},
			wantStatus:   http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mockLogic.NewMockController(ctrl)
			if tc.mockController != nil {
				tc.mockController(mock)
			}

			reprt := ReprtClient{
				deckController: mock,
				logger:         zerolog.Nop(),
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodGet, "/page/home", nil)
			require.NoError(t, err)

			req = req.WithContext(reptrCtx.AddUsername(req.Context(), tc.wantUserName))

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.HomePage(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			snaps.MatchSnapshot(t, rr.Body.String())
		})
	}
}

func TestReprtClient_CreateGroup(t *testing.T) {
	var (
		haveGroupID = "groupID"
	)
	testCases := map[string]struct {
		mockController func(mock *mockLogic.MockController)
		wantStatus     int
		haveUserName   string
		haveGroupName  string
	}{
		"should create group": {
			mockController: func(mock *mockLogic.MockController) {
				mock.EXPECT().CreateGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(haveGroupID, nil)
			},
			haveUserName:  "user",
			haveGroupName: "world",
			wantStatus:    http.StatusOK,
		},
		"should return internal error when username is not on context": {
			haveUserName: "",
			wantStatus:   http.StatusInternalServerError,
		},
		"should return 400 when group name is missing": {
			haveUserName:  "user",
			haveGroupName: "",
			wantStatus:    http.StatusBadRequest,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mockLogic.NewMockController(ctrl)
			if tc.mockController != nil {
				tc.mockController(mock)
			}

			reprt := ReprtClient{
				deckController: mock,
				logger:         zerolog.Nop(),
			}
			formValues := url.Values{}
			formValues.Add("group-name", tc.haveGroupName)

			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodPost, "/page/create-group", nil)
			require.NoError(t, err)

			req.PostForm = formValues

			req = req.WithContext(reptrCtx.AddUsername(req.Context(), tc.haveUserName))

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.CreateGroup(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			snaps.MatchSnapshot(t, rr.Body.String())
		})
	}
}

func TestReprtClient_CreateGroupPage(t *testing.T) {
	var ()
	testCases := map[string]struct {
		wantGroups   models.Group
		wantStatus   int
		wantUserName string
	}{
		"should load create group page": {
			wantUserName: "hello",
			wantStatus:   http.StatusOK,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {

			reprt := ReprtClient{
				logger: zerolog.Nop(),
			}
			// Create a request object with necessary parameters
			req, err := http.NewRequest(http.MethodGet, "/page/group", nil)
			require.NoError(t, err)

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the Login method with the mocked dependencies
			reprt.CreateGroupPage(rr, req)

			// Check the status code of the response
			assert.Equal(t, tc.wantStatus, rr.Code)
			snaps.MatchSnapshot(t, rr.Body.String())
		})
	}
}
