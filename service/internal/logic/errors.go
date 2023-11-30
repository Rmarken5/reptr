package logic

import "errors"

var (
	ErrInvalidToBeforeFrom = errors.New("'to' cannot be before 'from'")
	ErrInvalidGroupName    = errors.New("invalid group name")
)
