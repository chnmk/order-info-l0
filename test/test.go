package test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

// Удали меня
var E models.Order

// Удали меня
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

// Данные не претендуют на реалистичность, совпадение цен и так далее, поскольку по ТЗ мы пока ничего не делаем с этими данными, важно только то, что они есть.
func GenerateFakeData() []kafka.Message {

	var result []kafka.Message

	for i := 0; i < 1000; i++ {
		// В модель надо добавить теги вида `fake:"{number:1,100}"`
		// https://github.com/brianvoe/gofakeit
		var order models.Order

		// Желательно написать функцию для возвращения СТРОКИ в таком формате: 2021-11-26T06:22:19Z
		/*
			func (c *order.date_created) Fake(f *gofakeit.Faker) (any, error) {
				return f.DateRange(time.Now().AddDate(-100, 0, 0), time.Now().AddDate(-18, 0, 0)), nil
			}
		*/

		gofakeit.Struct(&order)

		// Посмотреть как она будет вести себя с []Item

		// Использовать marshall (?)

		// result = append(result, kafka.Message{Value: order})
	}

	return result
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

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

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

}
