package provider

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
	ErrNoUser     = errors.New("user does not exist")
)
