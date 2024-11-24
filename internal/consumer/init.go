package consumer

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func Init() {
	var conn *kafka.Conn
	var err error
	for i := 0; i < 20; i++ {
		conn, err = kafka.DialLeader(context.Background(), "tcp", "kafka:9092", "orders-1", 0)
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
