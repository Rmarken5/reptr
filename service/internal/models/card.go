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

	FrontOfCard struct {
		DeckID       string `bson:"deck_id"`
		CardID       string `bson:"card_id"`
		Content      string `bson:"content"`
		PreviousCard string `bson:"previous_card"`
		NextCard     string `bson:"next_card"`
		Upvotes      int    `bson:"upvotes"`
		Downvotes    int    `bson:"downvotes"`
	}

	IsUpvotedByUser   bool
	isDownvotedByUser bool

	BackOfCard struct {
		DeckID            string            `bson:"deck_id"`
		CardID            string            `bson:"card_id"`
		Answer            string            `bson:"answer"`
		NextCard          string            `bson:"next_card"`
		IsUpvotedByUser   IsUpvotedByUser   `bson:"is_upvoted_by_user"`
		IsDownvotedByUser isDownvotedByUser `bson:"is_downvoted_by_user"`
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

func (i IsUpvotedByUser) UpvotedClass() string {
	if i {
		return "upvoted"
	}
	return ""
}
func (i IsUpvotedByUser) NextUpvoteDirection() string {
	if i {
		return "remove_upvote"
	}
	return "upvote"
}

func (i isDownvotedByUser) DownvotedClass() string {
	if i {
		return "downvoted"
	}
	return ""
}
func (i IsUpvotedByUser) NextVoteDirection() string {
	if i {
		return "remove_Downvote"
	}
	return "downvote"
}
