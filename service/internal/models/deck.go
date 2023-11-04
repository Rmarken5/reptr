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

	WithCards struct {
		Deck
		Cards []Card
	}
)
