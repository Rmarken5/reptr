package models

import "time"

type (
	Group struct {
		ID         string     `bson:"_id"`
		Name       string     `bson:"name"`
		CreatedBy  string     `bson:"created_by"`
		Moderators []string   `bson:"moderators"`
		DeckIDs    []string   `bson:"deck_ids"`
		Members    []string   `bson:"members"`
		CreatedAt  time.Time  `bson:"created_at"`
		UpdatedAt  time.Time  `bson:"updated_at"`
		DeletedAt  *time.Time `bson:"deleted_at"`
	}

	HomePageGroup struct {
		ID         string     `bson:"_id"`
		Name       string     `bson:"name"`
		CreatedBy  string     `bson:"created_by"`
		Moderators []string   `bson:"moderators"`
		DeckIDs    []string   `bson:"deck_ids"`
		NumMembers int        `bson:"numMembers"`
		CreatedAt  time.Time  `bson:"created_at"`
		UpdatedAt  time.Time  `bson:"updated_at"`
		DeletedAt  *time.Time `bson:"deleted_at"`
	}

	GroupWithDecks struct {
		Group `bson:",inline"`
		Decks []GetDeckResults `bson:"decks"`
	}
)
