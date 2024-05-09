package database

import "errors"

var (
	ErrInsert = errors.New("error inserting")
	ErrUpdate = errors.New("error updating")
	ErrDelete = errors.New("error deleting")
	ErrFind   = errors.New("error finding")

	ErrAggregate = errors.New("error using aggregation")
	ErrNoResults = errors.New("results not found")
)
