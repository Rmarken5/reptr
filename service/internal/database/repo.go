package database

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination ./mocks/repo_mock.go -package database . Repository

var _ Repository = new(DAO)

type (
	Repository interface {
		CardDataAccess
		DeckDataAccess
		GroupDataAccess
		ProviderUsersDataAccess
		UserDataAccess
		SessionDataAccess
	}
	DAO struct {
		CardDataAccess
		DeckDataAccess
		GroupDataAccess
		ProviderUsersDataAccess
		UserDataAccess
		SessionDataAccess
	}
)

func NewRepository(logger zerolog.Logger, db *mongo.Database) *DAO {
	l := logger.With().Str("module", "Repository").Logger()
	return &DAO{
		NewCardDataAccess(db, l),
		NewDeckDataAccess(db, l),
		NewGroupDataAccess(db, l),
		NewProviderUsersDataAccess(db, l),
		NewUserDataAccess(db, l),
		NewSessionDataAccess(db, l),
	}
}
