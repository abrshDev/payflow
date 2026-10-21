package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokerAddr, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:                []string{brokerAddr},
			Topic:                  topic,
			GroupID:                groupID,
			WatchPartitionChanges:  true,
			PartitionWatchInterval: 5 * time.Second,
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler func(ctx context.Context, eventID string, payload []byte) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // context cancelled — this is a clean shutdown, not a real error
			}
			log.Printf("consumer: read failed: %v", err)
			continue
		}

		eventID := headerValue(msg.Headers, "event_id")

		if err := handler(ctx, eventID, msg.Value); err != nil {
			log.Printf("consumer: handler failed for event %s: %v", eventID, err)
		}
	}
}

func headerValue(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
