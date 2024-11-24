package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/segmentio/kafka-go"
)

func Init() {
	var conn *kafka.Conn
	var err error

	att, err := strconv.Atoi(cfg.Env["KAFKA_RECONNECT_ATTEMPTS"])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	for i := 0; i < att; i++ {
		conn, err = kafka.DialLeader(context.Background(),
			cfg.Env["KAFKA_NETWORK"],
			fmt.Sprintf("%s:%s", cfg.Env["KAFKA_PROTOCOL"], cfg.Env["KAFKA_PORT"]),
			cfg.Env["KAFKA_TOPIC"],
			0)

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
