package commands

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/abrshDev/ledger-service/internal/domain/entities"
	domainerrors "github.com/abrshDev/ledger-service/internal/domain/errors"
	"github.com/abrshDev/ledger-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type paymentEventPayload struct {
	PaymentID  uuid.UUID `json:"payment_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Amount     int64     `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
}

type ProcessPaymentEventHandler struct {
	repo repositories.LedgerRepository
}

func NewProcessPaymentEventHandler(repo repositories.LedgerRepository) *ProcessPaymentEventHandler {
	return &ProcessPaymentEventHandler{repo: repo}
}

func (h *ProcessPaymentEventHandler) Handle(ctx context.Context, eventID string, payload []byte) error {
	var event paymentEventPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	var direction entities.Direction
	switch event.Status {
	case "captured":
		direction = entities.DirectionCredit
	case "refunded":
		direction = entities.DirectionDebit
	default:
		return nil
	}

	entry := entities.NewLedgerEntry(event.PaymentID, event.Amount, event.Currency, direction, eventID)

	err := h.repo.Save(ctx, entry)
	if errors.Is(err, domainerrors.ErrDuplicateEvent) {
		return nil
	}
	return err
}