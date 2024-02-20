package middlewares

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/service/internal/api"
	mocks "github.com/rmarken/reptr/service/internal/api/middlewares/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	var (
		haveToken = "hello"
	)

	testCases := map[string]struct {
		mockStore  func(mock *mocks.MockStore)
		wantStatus int
		wantToken  string
	}{
		"should set token on request header when available": {
			mockStore: func(mock *mocks.MockStore) {
				mock.EXPECT().Get(gomock.Any(), api.CookieSessionID).Return(&sessions.Session{
					ID: uuid.NewString(),
					Values: map[interface{}]interface{}{
						api.SessionTokenKey: haveToken,
					},
					Options: nil,
					IsNew:   false,
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantToken:  haveToken,
		},
		"should return bad request when token is not set on session": {
			mockStore: func(mock *mocks.MockStore) {
				mock.EXPECT().Get(gomock.Any(), api.CookieSessionID).Return(&sessions.Session{
					ID:      uuid.NewString(),
					Values:  map[interface{}]interface{}{},
					Options: nil,
					IsNew:   false,
				}, nil)
			},
			wantStatus: http.StatusBadRequest,
			wantToken:  "",
		},
		"should return bad request when token is not a string": {
			mockStore: func(mock *mocks.MockStore) {
				mock.EXPECT().Get(gomock.Any(), api.CookieSessionID).Return(&sessions.Session{
					ID: uuid.NewString(),
					Values: map[interface{}]interface{}{
						api.SessionTokenKey: 123,
					},
					Options: nil,
					IsNew:   false,
				}, nil)
			},
			wantStatus: http.StatusBadRequest,
			wantToken:  "",
		},
		"should return Okay with no auth header when session is new": {
			mockStore: func(mock *mocks.MockStore) {
				mock.EXPECT().Get(gomock.Any(), api.CookieSessionID).Return(&sessions.Session{
					ID:      uuid.NewString(),
					Values:  map[interface{}]interface{}{},
					Options: nil,
					IsNew:   true,
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantToken:  "",
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockStore(ctrl)

			if tc.mockStore != nil {
				tc.mockStore(mock)
			}

			var gotToken string
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("called"))
				gotToken = request.Header.Get("Authorization")
			})
			sessionHandler := Session(zerolog.Nop(), mock)(handler)

			ts := httptest.NewServer(sessionHandler)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			_, err = io.ReadAll(res.Body)
			require.NoError(t, err)
			res.Body.Close()

			assert.Equal(t, tc.wantToken, gotToken)
		})
	}
}
