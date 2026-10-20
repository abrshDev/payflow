package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/abrshDev/ledger-service/internal/domain/entities"
	domainerrors "github.com/abrshDev/ledger-service/internal/domain/errors"
	"github.com/lib/pq"
)

type PostgresLedgerRepository struct {
	db *sql.DB
}

func NewPostgresLedgerRepository(db *sql.DB) *PostgresLedgerRepository {
	return &PostgresLedgerRepository{db: db}
}

func (r *PostgresLedgerRepository) Save(ctx context.Context, entry *entities.LedgerEntry) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ledger_entries (id, payment_id, amount, currency, direction, source_event_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, entry.ID, entry.PaymentID, entry.Amount, entry.Currency, entry.Direction, entry.SourceEventID, entry.CreatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domainerrors.ErrDuplicateEvent
		}
		return err
	}

	return nil
}
