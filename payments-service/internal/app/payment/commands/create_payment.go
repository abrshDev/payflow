package commands

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/entities"
	"github.com/abrshDev/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type CreatePaymentCommand struct {
	MerchantID uuid.UUID
	Amount     int64
	Currency   string
}

type CreatePaymentHandler struct {
	repo repositories.PaymentRepository
}

func NewCreatePaymentHandler(repo repositories.PaymentRepository) *CreatePaymentHandler {
	return &CreatePaymentHandler{repo: repo}
}

func (h *CreatePaymentHandler) Handle(ctx context.Context, cmd CreatePaymentCommand) (*entities.Payment, error) {
	payment, err := entities.NewPayment(cmd.MerchantID, cmd.Amount, cmd.Currency)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
