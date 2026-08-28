package entities

import (
	"time"

	domainerrors "github.com/abrshDev/payment-service/internal/domain/errors"
	"github.com/google/uuid"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusAuthorized Status = "authorized"
	StatusCaptured   Status = "captured"
	StatusFailed     Status = "failed"
	StatusRefunded   Status = "refunded"
)

type StatusTransition struct {
	From   Status
	To     Status
	At     time.Time
	Reason string
}

type Payment struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
	Amount     int64  // stored in the smallest currency unit (cents) — never float64 for money
	Currency   string // ISO 4217, e.g. "USD"
	Status     Status
	CreatedAt  time.Time
	History    []StatusTransition
}

func NewPayment(merchantID uuid.UUID, amount int64, currency string) (*Payment, error) {
	if amount <= 0 {
		return nil, domainerrors.ErrInvalidAmount
	}
	now := time.Now().UTC()
	return &Payment{
		ID:         uuid.New(),
		MerchantID: merchantID,
		Amount:     amount,
		Currency:   currency,
		Status:     StatusCreated,
		CreatedAt:  now,
		History: []StatusTransition{
			{From: "", To: StatusCreated, At: now, Reason: "payment created"},
		},
	}, nil
}

func (p *Payment) transition(to Status, reason string) {
	p.History = append(p.History, StatusTransition{
		From:   p.Status,
		To:     to,
		At:     time.Now().UTC(),
		Reason: reason,
	})
	p.Status = to
}

func (p *Payment) Authorize() error {
	if p.Status != StatusCreated {
		return domainerrors.ErrInvalidStateTransition
	}
	p.transition(StatusAuthorized, "payment authorized")
	return nil
}

func (p *Payment) Capture() error {
	if p.Status != StatusAuthorized {
		return domainerrors.ErrPaymentNotAuthorized
	}
	p.transition(StatusCaptured, "payment captured")
	return nil
}

func (p *Payment) Fail(reason string) error {
	if p.Status == StatusCaptured || p.Status == StatusFailed || p.Status == StatusRefunded {
		return domainerrors.ErrPaymentAlreadyTerminal
	}
	p.transition(StatusFailed, reason)
	return nil
}

func (p *Payment) Refund() error {
	if p.Status != StatusCaptured {
		return domainerrors.ErrPaymentNotCaptured
	}
	p.transition(StatusRefunded, "payment refunded")
	return nil
}
