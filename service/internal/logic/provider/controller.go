package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
)

var _ Logic = Controller{}

// Logic represents an interface for the business logic operations.
type (
	Logic interface {
		Authenticate(ctx context.Context, username, password string) (models.TokenResp, error)
		UserExists(ctx context.Context, subject string) (string, bool, error)
		InsertPair(ctx context.Context, subject string) error
	}

	Controller struct {
		clientID     string
		clientSecret string
		httpClient   http.Client
		authEndpoint string
		logger       zerolog.Logger
		repo         database.Repository
	}
)

// New initializes a new instance of the `Controller` struct with the provided logger and repository. It returns a pointer to the newly created `Controller`.
//
// Parameters:
// - logger: The logger instance used for logging within the controller.
// - repo: The database repository implementation.
//
// Returns:
// - *Controller: A pointer to the newly created `Controller` instance.
func New(logger zerolog.Logger, clientID, clientSecret, authEndpoint string, httpClient http.Client, repo database.Repository) *Controller {
	log := logger.With().Str("module", "provider logic").Logger()
	return &Controller{
		clientID:     clientID,
		clientSecret: clientSecret,
		authEndpoint: authEndpoint,
		httpClient:   httpClient,
		logger:       log,
		repo:         repo,
	}
}

func (c *Controller) Authenticate(ctx context.Context, username, password string) (models.TokenResp, error) {
	logger := c.logger.With().Str("model", "authenticate").Logger()
	logger.Info().Msgf("authenticating: %s", username)

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", username)
	form.Set("password", password)
	form.Set("audience", "reptr")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)

	payload := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest(http.MethodPost, c.authEndpoint, payload)
	if err != nil {
		logger.Error().Err(err).Msgf("error creating http request for authentication: %s", username)
		return models.TokenResp{}, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Msgf("while trying to authenticate: %s", username)
		return models.TokenResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Msgf("authentication failed for user: %s", username)
		return models.TokenResp{}, errors.New("authentication failed")
	}

	var tokenResp models.TokenResp
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		logger.Error().Err(err).Msgf("error decoding authentication response for user: %s", username)
		return models.TokenResp{}, err
	}

	logger.Info().Msgf("authentication successful for user: %s", username)
	return tokenResp, nil
}

// UserExists checks if a user exists for the given subject.
// It queries the database to get the user ID for the subject.
// If the user exists, it returns the user ID and true. Otherwise, it returns an empty string and false.
// If there is an error while querying the database, it logs the error and returns an empty string and false.
func (c *Controller) UserExists(ctx context.Context, subject string) (string, bool, error) {
	logger := c.logger.With().Str("method", "userExists").Logger()
	logger.Info().Msgf("checking if user exists for %s", subject)

	userID, err := c.repo.GetUserIDFor(ctx, subject)
	if err != nil {
		if errors.Is(err, database.ErrNoResults) {
			return "", false, nil
		}
		logger.Error().Err(err).Msgf("error getting user ID for subject %s", subject)
		return "", false, err
	}
	return userID, true, nil
}

// InsertPair inserts a user/subject pair into the database.
// It first checks if the user already exists. If the user exists, it returns an error.
// Then it inserts the user/subject pair into the database using a randomly generated UUID for the user ID.
// If there is an error during the insertion, it returns the error.
// Otherwise, it returns nil.
func (c *Controller) InsertPair(ctx context.Context, subject string) error {
	logger := c.logger.With().Str("method", "insertPair").Logger()
	logger.Info().Msgf("checking if user exists for %s", subject)

	_, userExists, err := c.UserExists(ctx, subject)
	if err != nil {
		logger.Error().Err(err).Msgf("error checking if user exists for subject %s", subject)
		return err
	}
	if userExists {
		logger.Debug().Msgf("user %s already exists", subject)
		return ErrUserExists
	}

	err = c.repo.InsertUserSubjectPair(ctx, uuid.NewString(), subject)
	if err != nil {
		logger.Error().Err(err).Msgf("error inserting user/subject pair for subject %s", subject)
		return err
	}
	return nil
}
