package commands

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type CapturePaymentCommand struct {
	PaymentID uuid.UUID
}

type CapturePaymentHandler struct {
	repo repositories.PaymentRepository
}

func NewCapturePaymentHandler(repo repositories.PaymentRepository) *CapturePaymentHandler {
	return &CapturePaymentHandler{
		repo: repo,
	}
}

func (h *CapturePaymentHandler) Handle(ctx context.Context, cmd CapturePaymentCommand) error {
	payment, err := h.repo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return err
	}

	if err := payment.Capture(); err != nil {
		return err
	}

	return h.repo.Save(ctx, payment)
}
