package mongorepo

import "errors"

var (
	errInvalidNumber = errors.New("number cannot be negative or be 0")
)
