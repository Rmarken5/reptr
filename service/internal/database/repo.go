package database

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination ./mocks/repo_mock.go -package database . Repository

var _ Repository = new(DataAccessObject)

type (
	Repository interface {
		CardDataAccess
		DeckDataAccess
		GroupDataAccess
		ProviderUsersDataAccess
		UserDataAccess
		SessionDataAccess
	}
	DataAccessObject struct {
		*CardDAO
		*DeckDAO
		*GroupDAO
		*ProviderUsersDAO
		*UserDAO
		*SessionDAO
	}
)

func NewRepository(logger zerolog.Logger, db *mongo.Database) *DataAccessObject {
	l := logger.With().Str("module", "Repository").Logger()
	return &DataAccessObject{
		NewCardDataAccess(db, l),
		NewDeckDataAccess(db, l),
		NewGroupDataAccess(db, l),
		NewProviderUsersDataAccess(db, l),
		NewUserDataAccess(db, l),
		NewSessionDataAccess(db, l),
	}
}
