package test

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

func ReadModelFile() (models.Order, []byte) {
	var E models.Order
	content, err := os.ReadFile("test/model.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &E)
	if err != nil {
		log.Fatal(err)
	}

	return E, content
}

// Данные не претендуют на реалистичность, совпадение цен и так далее, поскольку по ТЗ мы пока ничего не делаем с этими данными, важно только то, что они есть.
func GenerateFakeData() []kafka.Message {
	var result []kafka.Message

	slog.Info("generating fake data...")

	for i := 0; i < 10; i++ {
		var order models.Order

		// Желательно написать функцию для возвращения СТРОКИ в таком формате: 2021-11-26T06:22:19Z
		/*
			func (c *order.date_created) Fake(f *gofakeit.Faker) (any, error) {
				return f.DateRange(time.Now().AddDate(-100, 0, 0), time.Now().AddDate(-18, 0, 0)), nil
			}
		*/

		gofakeit.Struct(&order)

		// Посмотреть как она будет вести себя с []Item

		msg, err := json.Marshal(order)
		if err != nil {
			slog.Error(err.Error())
		}

		result = append(result, kafka.Message{Value: msg})
	}

	slog.Info("fake data generation finished")

	return result
}

func PublishTestData() {
	topic := "orders"
	partition := 0

	slog.Info("connecting to kafka to publish test data...")

	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	_, content := ReadModelFile()

	slog.Info("writing invalid data...")

	// Записывает некорректные данные
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

	slog.Info("writing fake data...")

	// Записывает моковые данные
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		GenerateFakeData()...,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	slog.Info("test data successfully published")

}
