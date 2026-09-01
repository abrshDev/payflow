package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/abrshDev/payment-service/internal/domain/entities"
	"github.com/google/uuid"
)

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Save(ctx context.Context, payment *entities.Payment) error {
	historyJSON, err := json.Marshal(payment.History)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO payments (id, merchant_id, amount, currency, status, created_at, history)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET status = EXCLUDED.status, history = EXCLUDED.history
	`

	_, err = r.db.ExecContext(ctx, query,
		payment.ID, payment.MerchantID, payment.Amount, payment.Currency,
		payment.Status, payment.CreatedAt, historyJSON,
	)
	return err
}

func (r *PostgresPaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, created_at, history
		FROM payments
		WHERE id = $1
	`

	var p entities.Payment
	var historyJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.MerchantID, &p.Amount, &p.Currency,
		&p.Status, &p.CreatedAt, &historyJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(historyJSON, &p.History); err != nil {
		return nil, err
	}

	return &p, nil
}
