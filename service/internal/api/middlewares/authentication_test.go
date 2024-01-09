package middlewares

import (
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	auth "github.com/rmarken/reptr/service/internal/logic/auth/mocks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {

	testCases := map[string]struct {
		wantToken    string
		wantResponse string
		wantStatus   int
		mockAuth     func(mock *auth.MockAuthentication)
		wantErr      error
		wantUserID   string
	}{
		"should serve endpoint": {
			wantToken:    "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHwxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.-dxzIIdM5cHR_yNnAtkxtUIIZhjLkHOeMEoUGurb_ho",
			wantResponse: "called",
			wantStatus:   http.StatusOK,
			mockAuth: func(mock *auth.MockAuthentication) {
				mock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{
					Subject: "auth0|1234567890",
				}, nil)
			},
			wantUserID: "1234567890",
		},
		"should return status forbidden": {
			wantToken:    "",
			wantResponse: "",
			wantStatus:   http.StatusForbidden,
		},
		"should return status unauthorized": {
			wantToken:    "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHwxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.-dxzIIdM5cHR_yNnAtkxtUIIZhjLkHOeMEoUGurb_ho",
			wantResponse: "",
			mockAuth: func(mock *auth.MockAuthentication) {
				mock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(nil, errors.New("not authorized"))
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockAuth := auth.NewMockAuthentication(ctrl)
			if tc.mockAuth != nil {
				tc.mockAuth(mockAuth)
			}
			var userID string
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("called"))
				userID = request.Context().Value(models.UserIDKey).(string)
			})

			authHandler := Authenticate(zerolog.Nop(), mockAuth)(handler)
			ts := httptest.NewServer(authHandler)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", tc.wantToken)

			res, err := ts.Client().Do(req)

			require.ErrorIs(t, err, tc.wantErr)

			resp, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			res.Body.Close()

			assert.Equal(t, tc.wantResponse, string(resp))
			assert.Equal(t, tc.wantUserID, userID)
		})

	}
}
