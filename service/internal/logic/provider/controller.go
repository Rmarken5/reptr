package provider

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rs/zerolog"
)

var _ Controller = &Logic{}

// Logic represents an interface for the business logic operations.
type (
	Controller interface {
		UserExists(ctx context.Context, subject string) (string, bool, error)
		InsertPair(ctx context.Context, subject string) error
	}

	Logic struct {
		logger zerolog.Logger
		repo   database.Repository
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
func New(logger zerolog.Logger, repo database.Repository) *Logic {
	log := logger.With().Str("module", "provider logic").Logger()
	return &Logic{
		logger: log,
		repo:   repo,
	}
}

// UserExists checks if a user exists for the given subject.
// It queries the database to get the user ID for the subject.
// If the user exists, it returns the user ID and true. Otherwise, it returns an empty string and false.
// If there is an error while querying the database, it logs the error and returns an empty string and false.
func (l *Logic) UserExists(ctx context.Context, subject string) (string, bool, error) {
	logger := l.logger.With().Str("method", "userExists").Logger()
	logger.Info().Msgf("checking if user exists for %s", subject)

	userID, err := l.repo.GetUserIDFor(ctx, subject)
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
func (l *Logic) InsertPair(ctx context.Context, subject string) error {
	logger := l.logger.With().Str("method", "insertPair").Logger()
	logger.Info().Msgf("checking if user exists for %s", subject)

	_, userExists, err := l.UserExists(ctx, subject)
	if err != nil {
		logger.Error().Err(err).Msgf("error checking if user exists for subject %s", subject)
		return err
	}
	if userExists {
		logger.Debug().Msgf("user %s already exists", subject)
		return ErrUserExists
	}

	err = l.repo.InsertUserSubjectPair(ctx, uuid.NewString(), subject)
	if err != nil {
		logger.Error().Err(err).Msgf("error inserting user/subject pair for subject %s", subject)
		return err
	}
	return nil
}
