package models

import (
	"time"
)

type (
	Deck struct {
		ID           string    `bson:"_id"`
		Name         string    `bson:"name"`
		UserUpvote   []string  `bson:"user_upvote"`
		UserDownvote []string  `bson:"user_downvote"`
		CreatedAt    time.Time `bson:"created_at"`
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
)
