package models

import (
	"time"
)

type (
	Deck struct {
		ID        string    `bson:"_id"`
		Name      string    `bson:"name"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}

	DeckWithCards struct {
		Deck  `bson:",inline"`
		Cards []Card
	}
)
