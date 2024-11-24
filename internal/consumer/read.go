package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/segmentio/kafka-go"
)

/*
TODO: написать объяснительную.

Consumer group используются по нескольким причинам:
	- Возможность напрямую использовать коммиты и не перечитывать старые сообщения.
*/

func (c *KafkaConsumer) Read(ctx context.Context) {
	slog.Info("creating new kafka reader...")

	maxBytes, err := strconv.Atoi(cfg.Env["KAFKA_MAX_BYTES"])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	att, err := strconv.Atoi(cfg.Env["KAFKA_RECONNECT_ATTEMPTS"])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	interv, err := strconv.Atoi(cfg.Env["KAFKA_COMMIT_INVERVAL_SECONDS"])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{fmt.Sprintf("%s:%s", cfg.Env["KAFKA_PROTOCOL"], cfg.Env["KAFKA_PORT"])},
		GroupID:        cfg.Env["KAFKA_GROUP_ID"],
		Topic:          cfg.Env["KAFKA_TOPIC"],
		MaxBytes:       maxBytes,
		MaxAttempts:    att,
		CommitInterval: time.Duration(interv) * time.Second,
	})

	slog.Info("reader created, reding messages...")

	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		slog.Info("=== handling new order ===")
		memory.DATA.HandleMessage(m.Value)
	}

	if err := c.Reader.Close(); err != nil {
		slog.Error("failed to close reader: " + err.Error())
	}

	slog.Info("closing consumer connection...")
}
