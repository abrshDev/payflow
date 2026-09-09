package outbox

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/abrshDev/payment-service/internal/infrastructure/kafka"
)

type Relay struct {
	db       *sql.DB
	producer *kafka.Producer
}

func NewRelay(db *sql.DB, producer *kafka.Producer) *Relay {
	return &Relay{db: db, producer: producer}
}

func (r *Relay) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.publishPending(ctx)
		}
	}
}

func (r *Relay) publishPending(ctx context.Context) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, aggregate_id, event_type, payload
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY created_at
	`)
	if err != nil {
		log.Printf("relay: query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, aggregateID, eventType string
		var payload []byte
		if err := rows.Scan(&id, &aggregateID, &eventType, &payload); err != nil {
			log.Printf("relay: scan failed: %v", err)
			continue
		}

		if err := r.producer.Publish(ctx, aggregateID, payload); err != nil {
			log.Printf("relay: publish failed for event %s: %v", id, err)
			continue
		}

		_, err := r.db.ExecContext(ctx, `UPDATE outbox_events SET published_at = $1 WHERE id = $2`, time.Now().UTC(), id)
		if err != nil {
			log.Printf("relay: failed to mark published for event %s: %v", id, err)
		}
	}
}
