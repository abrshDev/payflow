package repositories

import (
	"context"

	"github.com/abrshDev/payment-service/internal/domain/entities"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *entities.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
}
