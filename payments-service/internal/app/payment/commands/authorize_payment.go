package commands

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type AuthorizePaymentCommand struct {
	PaymentID uuid.UUID
}

type AuthorizePaymentHandler struct {
	repo repositories.PaymentRepository
}

func NewAuthorizePaymentHandler(repo repositories.PaymentRepository) *AuthorizePaymentHandler {
	return &AuthorizePaymentHandler{repo: repo}
}

func (h *AuthorizePaymentHandler) Handle(ctx context.Context, cmd AuthorizePaymentCommand) error {
	payment, err := h.repo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return err
	}

	if err := payment.Authorize(); err != nil {
		return err
	}

	return h.repo.Save(ctx, payment)
}
