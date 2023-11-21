package models

import (
	"time"
)

type (
	Deck struct {
		ID        string    `bson:"_id"`
		Name      string    `bson:"name"`
		CreatedAt time.Time `bson:"created_at"`
	}

	DeckWithCards struct {
		Deck  `bson:",inline"`
		Cards []Card
	}
)
