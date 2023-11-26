package models

import "time"

type (
	Group struct {
		ID        string     `bson:"_id"`
		Name      string     `bson:"name"`
		CreatedAt time.Time  `bson:"created_at"`
		UpdatedAt time.Time  `bson:"updated_at"`
		DeletedAt *time.Time `bson:"deleted_at"`
	}

	GroupWithDecks struct {
		Group `bson:",inline"`
		Decks []Deck
	}
)
