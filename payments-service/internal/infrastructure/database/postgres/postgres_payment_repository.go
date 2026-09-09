package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	historyJSON, err := json.Marshal(payment.History)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO payments (id, merchant_id, amount, currency, status, created_at, history)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET status = EXCLUDED.status, history = EXCLUDED.history
	`, payment.ID, payment.MerchantID, payment.Amount, payment.Currency,
		payment.Status, payment.CreatedAt, historyJSON)
	if err != nil {
		return err
	}

	for _, event := range payment.PullEvents() {
		payloadJSON, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO outbox_events (id, aggregate_id, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New(), payment.ID, event.Type, payloadJSON, time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
