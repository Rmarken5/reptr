package models

import "time"

type (
	Group struct {
		ID        string
		Name      string
		DeckIDs   []string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	GroupWithDecks struct {
		Group `bson:",inline"`
		Decks []Deck
	}
)
