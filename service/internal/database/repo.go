package database

import (
	"github.com/rmarken/reptr/internal/database/card"
	"github.com/rmarken/reptr/internal/database/deck"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	cardDAO card.DataAccess
	deckDAO deck.DataAccess
}

func NewRepository(logger zerolog.Logger, db *mongo.Database) *Repository {

}
