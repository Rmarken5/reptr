// platform/authenticator/auth.go

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

var _ Authentication = new(Authenticator)

//go:generate mockgen -destination ./mocks/controller_mock.go -package auth . Authentication
type (
	Authentication interface {
		VerifyIDToken(ctx context.Context, rawToken string) (*oidc.IDToken, error)
		PasswordCredentialsToken(ctx context.Context, username string, password string) (*oauth2.Token, error)
		RegisterUser(ctx context.Context, username, password string) (models.RegistrationUser, models.RegistrationError, error)
	}
	// Authenticator is used to authenticate our users.
	Authenticator struct {
		*oidc.Provider
		oauth2.Config
		audience   string
		endpoint   string
		logger     zerolog.Logger
		httpClient http.Client
	}
)

// New instantiates the *Authenticator.
func New(ctx context.Context, logger zerolog.Logger, audience, endpoint, clientID, clientSecret, callbackURL string) (*Authenticator, error) {
	log := logger.With().Str("module", "authenticator").Logger()
	provider, err := oidc.NewProvider(
		ctx,
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		endpoint: endpoint,
		audience: audience,
		Provider: provider,
		Config:   conf,
		logger:   log,
	}, nil
}

// VerifyIDToken verifies that an token
func (a *Authenticator) VerifyIDToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}
	return a.Verifier(oidcConfig).Verify(ctx, rawToken)
}
func (a *Authenticator) RegisterUser(ctx context.Context, username, password string) (models.RegistrationUser, models.RegistrationError, error) {
	log := a.logger.With().Str("module", "registerUser").Logger()

	reqBody := models.RegistrationRequest{
		ClientID:   a.Config.ClientID,
		Email:      username,
		Password:   password,
		Connection: "Username-Password-Authentication",
		Username:   username,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return models.RegistrationUser{}, models.RegistrationError{}, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sdbconnections/signup", a.endpoint), bytes.NewReader(jsonBytes))
	if err != nil {
		return models.RegistrationUser{}, models.RegistrationError{}, err
	}
	h := request.Header
	h.Set("Content-Type", "application/json")
	request.Header = h

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return models.RegistrationUser{}, models.RegistrationError{}, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.RegistrationUser{}, models.RegistrationError{}, err
	}
	log.Debug().Msgf("resp: %s", respBytes)

	var registerUser models.RegistrationUser
	err = json.Unmarshal(respBytes, &registerUser)
	if err != nil {
		return models.RegistrationUser{}, models.RegistrationError{}, err
	}

	if registerUser.IsZero() {
		var registrationErr models.RegistrationError
		err := json.Unmarshal(respBytes, &registrationErr)
		if err != nil {
			return models.RegistrationUser{}, models.RegistrationError{}, err
		}
		log.Debug().Msgf("reg err: %+v", registrationErr)
		return models.RegistrationUser{}, registrationErr, nil
	}

	return registerUser, models.RegistrationError{}, nil
}
