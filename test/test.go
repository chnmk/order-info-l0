package test

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishTestData() {
	topic := "orders"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("A: " + time.Now().String())},
		kafka.Message{Value: []byte("B: " + time.Now().String())},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
