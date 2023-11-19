package database

import "errors"

var (
	ErrInsert = errors.New("error inserting")
	ErrUpdate = errors.New("error updating")

	ErrAggregate = errors.New("error using aggregation")
	ErrNoResults = errors.New("results not found")
)
