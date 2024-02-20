package models

import (
	"time"
)

type (
	Deck struct {
		ID           string    `bson:"_id"`
		Name         string    `bson:"name"`
		UserUpvote   []string  `bson:"user_upvotes"`
		UserDownvote []string  `bson:"user_downvotes"`
		CreatedAt    time.Time `bson:"created_at"`
		CreatedBy    string    `bson:"created_by"`
		UpdatedAt    time.Time `bson:"updated_at"`
	}

	GetDeckResults struct {
		ID        string    `bson:"_id"`
		Name      string    `bson:"name"`
		Upvotes   int       `bson:"upvotes"`
		Downvotes int       `bson:"downvotes"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
	DeckWithCards struct {
		GetDeckResults `bson:",inline"`
		Cards          []Card
	}

	DeckSession struct {
		ID            string `bson:"_id"`
		DeckName      string `bson:"deck_name"`
		CurrentCardID string `bson:"current_card_id"`
		IsFront       bool   `bson:"is_front"`
	}
)
