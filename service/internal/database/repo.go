package database

import (
	"github.com/rmarken/reptr/internal/database/card"
	"github.com/rmarken/reptr/internal/database/deck"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Repository interface {
		card.CardDataAccess
		deck.DeckDataAccess
	}
	DAO struct {
		card.CardDataAccess
		deck.DeckDataAccess
	}
)

func NewRepository(logger zerolog.Logger, db *mongo.Database) *DAO {
	l := logger.With().Str("module", "Repository").Logger()
	return &DAO{
		card.NewDataAccess(db, l),
		deck.NewDataAccess(db, l),
	}
}
