package test

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"math/rand"
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

/*
Генерирует случайные корректные данные с использованием Gofakeit.

Данные не претендуют на реалистичность (не совпадают цены, места и так далее), но это и не мусорные рандомные символы.
Такой вариант приемлем, поскольку по ТЗ эти данные пока никак не обрабатываются, но отобразить их нужно.
*/
func GenerateFakeData() []kafka.Message {
	var result []kafka.Message

	slog.Info("generating fake data...")

	// Пример кастомной функции для генерации приближенных к реальности данных.
	gofakeit.AddFuncLookup("wbdate", gofakeit.Info{
		Category:    "custom",
		Description: "random date string",
		Example:     "2021-11-26T06:22:19Z",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			// Отнимает один мечяц от сегодняшнего дня, переводит в формат unix.
			min := time.Now().AddDate(0, -1, 0).Unix()
			// Прибавляет к min случайное значение, таким образом получает дату между сегодня и месяц назад.
			unix := min + rand.Int63n(time.Now().Unix()-min)

			time := time.Unix(unix, 0).Format("2006-01-02T15:04:05Z")
			return time, nil
		},
	})

	// Генерируем 10 сообщений для кафки (TODO: скорее всего будем генерировать по одному, а кастомную функцию отдельно вынесем)
	for i := 0; i < 10; i++ {
		var order models.Order
		gofakeit.Struct(&order)

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
