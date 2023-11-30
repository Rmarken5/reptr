package logic

import "errors"

var (
	ErrInvalidToBeforeFrom = errors.New("'to' cannot be before 'from'")
	ErrInvalidGroupName    = errors.New("invalid group name")
	ErrEmptyGroupID        = errors.New("empty group ID")
	ErrInvalidDeckName     = errors.New("invalid deck name")
	ErrEmptyDeckID         = errors.New("empty deck ID")
)
