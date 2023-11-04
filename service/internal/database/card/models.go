package card

import "time"

type (
	Type int
	Card struct {
		ID        string    `bson:"_id"`
		Front     string    `bson:"front"`
		Back      string    `bson:"back,omitempty"`
		Kind      Type      `bson:"type"`
		DeckID    string    `bson:"deck_id"`
		CreatedAt time.Time `bson:"created_at"`
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
