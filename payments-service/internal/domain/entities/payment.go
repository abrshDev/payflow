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
	From   Status    `json:"from"`
	To     Status    `json:"to"`
	At     time.Time `json:"at"`
	Reason string    `json:"reason"`
}

type Payment struct {
	ID         uuid.UUID          `json:"id"`
	MerchantID uuid.UUID          `json:"merchant_id"`
	Amount     int64              `json:"amount"`
	Currency   string             `json:"currency"`
	Status     Status             `json:"status"`
	CreatedAt  time.Time          `json:"created_at"`
	History    []StatusTransition `json:"history"`
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
