package dumb

type (
	CardFront struct {
		CardType       string
		CardID         string
		Front          string
		Upvotes        string
		Downvotes      string
		PreviousCardID string
		NextCardID     string
	}

	CardBack struct {
		CardID      string
		BackContent string
		NextCardID  string
		IsUpvoted   bool
		IsDownvoted bool
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
