package consumer

import (
	"context"
	"log/slog"
	"os"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/segmentio/kafka-go"
)

func Connect() {
	var conn *kafka.Conn
	var err error

	for i := 0; i < cfg.KafkaReconnectAttempts; i++ {
		conn, err = kafka.DialLeader(context.Background(),
			cfg.KafkaInitNework,
			cfg.KafkaInitAddress,
			cfg.KafkaInitTopic,
			cfg.KafkaInitPartition,
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

	if cfg.Env.Get("KAFKA_WRITE_EXAMPLES") == "1" {
		goFakeInit()

		for i := 0; i < cfg.KafkaWriterGoroutines; i++ {
			go publishExampleData()
		}

	}

	for i := 0; i < cfg.KafkaReaderGoroutines; i++ {
		go newConsumer().Read(context.TODO())
	}
}
