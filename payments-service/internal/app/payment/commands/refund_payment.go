package commands

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type RefundPaymentCommand struct {
	PaymentID uuid.UUID
}

type RefundPaymentHandler struct {
	repo repositories.PaymentRepository
}

func NewRefundPaymentHandler(repo repositories.PaymentRepository) *RefundPaymentHandler {
	return &RefundPaymentHandler{
		repo: repo,
	}
}

func (h *RefundPaymentHandler) Handle(ctx context.Context, cmd RefundPaymentCommand) error {

	payment, err := h.repo.FindByID(ctx, cmd.PaymentID)

	if err != nil {
		return err
	}
	if err := payment.Refund(); err != nil {
		return err
	}
	return h.repo.Save(ctx, payment)
}
