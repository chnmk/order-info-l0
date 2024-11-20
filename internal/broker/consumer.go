package broker

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consume() {
	topic := "orders"
	partition := 0

	var conn *kafka.Conn
	var err error

	for retry := 0; retry < 10; retry++ {
		conn, err = kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
		if err == nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	slog.Info("reading orders...")
	batch := conn.ReadBatch(1, 1e6) // fetch 1B min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
