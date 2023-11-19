package models

import "time"

type (
	Type int
	Card struct {
		ID        string    `bson:"_id"`
		Front     string    `bson:"front,omitempty"`
		Back      string    `bson:"back,omitempty"`
		Kind      Type      `bson:"type,omitempty"`
		DeckID    string    `bson:"deck_id,omitempty"`
		CreatedAt time.Time `bson:"created_at,omitempty"`
		UpdatedAt time.Time `bson:"update_at,omitempty"`
	}
)

const (
	BasicCard = iota
	MultipleChoice
)

func (c Type) String() string {
	switch c {
	case BasicCard:
		return "basic"
	case MultipleChoice:
		return "multiple choice"
	}
	return "unknown"
}
