package broker

import (
	"context"
	"log/slog"
	"time"

	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/segmentio/kafka-go"
)

/*
TODO: написать объяснительную.

Consumer group используются по нескольким причинам:
	- Возможность напрямую использовать коммиты и не перечитывать старые сообщения.
*/

func Consume(ctx context.Context) {
	slog.Info("creating new kafka reader...")

	// make a new reader that consumes from orders
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka:9092"},
		GroupID:        "go-orders-1",
		Topic:          "orders-1",
		MaxBytes:       100e3, // 100kb
		MaxAttempts:    20,
		CommitInterval: 1 * time.Second,
	})

	slog.Info("reader created, reding messages...")

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		slog.Info("=== handling new order ===")
		memory.DATA.HandleMessage(m.Value)
	}

	if err := r.Close(); err != nil {
		slog.Error("failed to close reader: " + err.Error())
	}

	slog.Info("closing consumer connection...")
}
