package models

import "strings"

type Vote int

const (
	Upvote Vote = iota
	Downvote
	RemoveUpvote
	RemoveDownvote
	Unknown
)

func (v Vote) String() string {
	switch v {
	case Upvote:
		return "upvote"
	case Downvote:
		return "downvote"
	case RemoveUpvote:
		return "remove_upvote"
	case RemoveDownvote:
		return "remove_downvote"
	default:
		return "unknown"
	}
}

func (v Vote) NextUpvote() Vote {
	if v == Upvote {
		return RemoveUpvote
	}
	return Upvote
}

func (v Vote) NextDownvote() Vote {
	if v == Downvote {
		return RemoveDownvote
	}
	return Downvote
}

func (v Vote) DownvoteClass() string {
	if v == Downvote {
		return "downvoted"
	}
	return ""
}

func (v Vote) UpvoteClass() string {
	if v == Upvote {
		return "upvoted"
	}
	return ""
}

func VoteFromString(direction string) Vote {
	direction = strings.ToLower(direction)
	switch direction {
	case "upvote":
		return Upvote
	case "downvote":
		return Downvote
	case "remove_upvote":
		return RemoveUpvote
	case "remove_downvote":
		return RemoveDownvote
	default:
		return Unknown
	}
}
