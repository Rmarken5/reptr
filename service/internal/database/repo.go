package database

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination ./mocks/repo_mock.go -package database . Repository

type (
	Repository interface {
		CardDataAccess
		DeckDataAccess
	}
	DAO struct {
		CardDataAccess
		DeckDataAccess
	}
)

func NewRepository(logger zerolog.Logger, db *mongo.Database) *DAO {
	l := logger.With().Str("module", "Repository").Logger()
	return &DAO{
		NewCardDataAccess(db, l),
		NewDeckDataAccess(db, l),
	}
}
