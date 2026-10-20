package entities

import (
	"time"

	"github.com/google/uuid"
)

type Direction string

const (
	DirectionCredit Direction = "credit"
	DirectionDebit  Direction = "debit"
)

type LedgerEntry struct {
	ID            uuid.UUID `json:"id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Direction     Direction `json:"direction"`
	SourceEventID string    `json:"source_event_id"` // the outbox event's own id — this is the idempotency key
	CreatedAt     time.Time `json:"created_at"`
}

func NewLedgerEntry(paymentID uuid.UUID, amount int64, currency string, direction Direction, sourceEventID string) *LedgerEntry {
	return &LedgerEntry{
		ID:            uuid.New(),
		PaymentID:     paymentID,
		Amount:        amount,
		Currency:      currency,
		Direction:     direction,
		SourceEventID: sourceEventID,
		CreatedAt:     time.Now().UTC(),
	}
}
