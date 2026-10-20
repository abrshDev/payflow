package repositories

import (
	"context"

	"github.com/abrshDev/ledger-service/internal/domain/entities"
)

type LedgerRepository interface {
	Save(ctx context.Context, entry *entities.LedgerEntry) error
}
