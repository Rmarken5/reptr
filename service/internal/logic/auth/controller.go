// platform/authenticator/auth.go

package auth

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

//go:generate mockgen -destination ./mocks/controller_mock.go -package auth . Authentication
type (
	Authentication interface {
		VerifyIDToken(ctx context.Context, rawToken string) (*oidc.IDToken, error)
		PasswordCredentialsToken(ctx context.Context, username string, password string) (*oauth2.Token, error)
	}
	// Authenticator is used to authenticate our users.
	Authenticator struct {
		*oidc.Provider
		oauth2.Config
	}
)

// New instantiates the *Authenticator.
func New(ctx context.Context, endpoint, clientID, clientSecret, callbackURL string) (*Authenticator, error) {
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
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an token
func (a *Authenticator) VerifyIDToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawToken)
}
