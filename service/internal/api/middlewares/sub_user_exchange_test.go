package middlewares

import (
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	auth "github.com/rmarken/reptr/service/internal/logic/auth/mocks"
	mocks "github.com/rmarken/reptr/service/internal/logic/provider/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeSubjectForUser(t *testing.T) {

	var (
		haveUserName = uuid.NewString()
		haveSubject  = uuid.NewString()
		token        = "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHwxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.-dxzIIdM5cHR_yNnAtkxtUIIZhjLkHOeMEoUGurb_ho"
	)

	testCases := map[string]struct {
		haveSubject  string
		mockProvider func(mock *mocks.MockController)
		mockAuth     func(mock *auth.MockAuthentication)
		wantStatus   int
		wantUserName string
	}{
		"should return username with 200": {
			mockProvider: func(mock *mocks.MockController) {
				mock.EXPECT().UserExists(gomock.Any(), haveSubject).Return(haveUserName, true, nil)
			},
			mockAuth: func(mock *auth.MockAuthentication) {
				mock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{Subject: haveSubject}, nil)
			},
			wantStatus:   http.StatusOK,
			wantUserName: haveUserName,
			haveSubject:  haveSubject,
		},
		"should return internal error when logic errors": {
			mockProvider: func(mock *mocks.MockController) {
				mock.EXPECT().UserExists(gomock.Any(), haveSubject).Return("", false, errors.New("idk"))
			},
			mockAuth: func(mock *auth.MockAuthentication) {
				mock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{Subject: haveSubject}, nil)
			},
			wantStatus:   http.StatusInternalServerError,
			wantUserName: "",
			haveSubject:  haveSubject,
		},
		"should return not found when username doesn't exist": {
			mockProvider: func(mock *mocks.MockController) {
				mock.EXPECT().UserExists(gomock.Any(), haveSubject).Return("", false, nil)
			},
			mockAuth: func(mock *auth.MockAuthentication) {
				mock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{Subject: haveSubject}, nil)
			},
			wantStatus:   http.StatusNotFound,
			wantUserName: "",
			haveSubject:  haveSubject,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockController(ctrl)
			mockAuth := auth.NewMockAuthentication(ctrl)

			if tc.mockProvider != nil {
				tc.mockProvider(mock)
			}

			if tc.mockAuth != nil {
				tc.mockAuth(mockAuth)
			}

			var gotUserName string
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("called"))
				gotUserName, _ = reptrCtx.Username(request.Context())
			})
			ex := ExchangeSubjectForUser(zerolog.Nop(), mock)(handler)
			a := Authenticate(zerolog.Nop(), mockAuth)
			h := a(ex)

			ts := httptest.NewServer(h)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", token)
			ctx := reptrCtx.AddSubject(req.Context(), tc.haveSubject)

			res, err := ts.Client().Do(req.WithContext(ctx))
			require.NoError(t, err)

			_, err = io.ReadAll(res.Body)
			require.NoError(t, err)
			res.Body.Close()

			assert.Equal(t, tc.wantUserName, gotUserName)

		})
	}
}
