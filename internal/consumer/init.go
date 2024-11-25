package consumer

import (
	"context"
	"log/slog"
	"os"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/segmentio/kafka-go"
)

func Init() {
	var conn *kafka.Conn
	var err error

	for i := 0; i < cfg.KafkaReconnectAttempts; i++ {
		conn, err = kafka.DialLeader(context.Background(),
			cfg.ConsumerInitNetwork,
			cfg.ConsumerInitAddress,
			cfg.ConsumerInitTopic,
			cfg.ConsumerInitPartition,
		)

		if err != nil {
			slog.Error(err.Error())
		} else {
			slog.Info("connection successful")
			break
		}

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		slog.Error("Kafka connection failed: " + err.Error())
		os.Exit(1)
	}

	conn.Close()
}
