package dumb

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
