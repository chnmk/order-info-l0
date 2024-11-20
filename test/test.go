package test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

// TEMP
var E models.Order

// TEMP
func ReadModelFile() {
	content, err := os.ReadFile("test/model.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &E)
	if err != nil {
		log.Fatal(err)
	}

}

func PublishTestData() {
	topic := "orders"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	var order models.Order
	content, err := os.ReadFile("test/model.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &order)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("Date: " + time.Now().String())},
		kafka.Message{Value: []byte("example")},
		kafka.Message{Value: []byte("\"{order_uid\": \"abcd\"")},
		kafka.Message{Value: []byte("\"{order_uid\": \"abcd\", \"example\": \"dcba\"}")},
		kafka.Message{Value: content},
		kafka.Message{Value: []byte("===============")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
