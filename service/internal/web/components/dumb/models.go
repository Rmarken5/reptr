package dumb

import "time"

type (
	CardFront struct {
		DeckID         string
		CardType       string
		CardID         string
		Front          string
		Upvotes        string
		Downvotes      string
		PreviousCardID string
		NextCardID     string
	}

	CardBack struct {
		DeckID         string
		CardID         string
		BackContent    string
		PreviousCardID string
		NextCardID     string
		IsUpvoted      bool
		IsDownvoted    bool
		VoteButtonData VoteButtonsData
	}

	// VoteButtonsData is data for the VoteButtons component
	VoteButtonsData struct {
		CardID            string
		UpvoteClass       string
		DownvoteClass     string
		UpvoteDirection   string
		DownvoteDirection string
	}
	Deck struct {
		ID           string
		DeckName     string
		NumUpvotes   int
		NumDownvotes int
		NumCards     int
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
)

func (c CardBack) UpvoteClass() string {
	if c.IsUpvoted {
		return "upvoted"
	}
	return ""
}

func (c CardBack) DownvoteClass() string {
	if c.IsDownvoted {
		return "downvoted"
	}
	return ""
}
