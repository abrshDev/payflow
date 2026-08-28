package errors

import "errors"

var (
	ErrInvalidStateTransition = errors.New("invalid payment state transition")
	ErrPaymentNotAuthorized   = errors.New("payment must be authorized before it can be captured")
	ErrPaymentNotCaptured     = errors.New("payment must be captured before it can be refunded")
	ErrInvalidAmount          = errors.New("payment amount must be greater than zero")
	ErrPaymentAlreadyTerminal = errors.New("payment is in a terminal state and cannot be modified")
)
