package errors

import "errors"

var (
	ErrDuplicateEvent = errors.New("event already processed")
)
