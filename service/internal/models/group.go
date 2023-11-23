package models

import "time"

type (
	Group struct {
		ID        string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	GroupWithDecks struct {
		Group `bson:",inline"`
		Decks []Deck
	}
)
