package queries

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/entities"
	"github.com/abrshDev/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type GetPaymentQuery struct {
	PaymentID uuid.UUID
}

type GetPaymentHandler struct {
	repo repositories.PaymentRepository
}

func NewGetPaymentHandler(repo repositories.PaymentRepository) *GetPaymentHandler {
	return &GetPaymentHandler{
		repo: repo,
	}
}

func (h *GetPaymentHandler) Handle(ctx context.Context, query GetPaymentQuery) (*entities.Payment, error) {
	payment, err := h.repo.FindByID(ctx, query.PaymentID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
