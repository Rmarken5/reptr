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
